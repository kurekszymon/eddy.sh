import { execFile } from 'child_process';
import fs from 'fs';
import path from 'path';
import https from 'https';

import { EDDY_BIN_DIR, EDDY_DIR } from '@/lib/consts';
import { logger } from '@/lib/logger';
import type { semver } from '@/lib/types';

export const ensureToolDir = (dirName: string) => {
    const dir = path.join(EDDY_DIR, dirName);

    if (fs.existsSync(dir)) {
        return dir;
    }

    try {
        fs.mkdirSync(dir, { mode: 0o755, recursive: true });
    } catch (e) {
        // TODO: handle error
        logger.error(e);
    }

    return dir;
};


export const downloadFile = (filePath: string, url: string, maxRedirects = 5): Promise<string> => {
    return new Promise((resolve, reject) => {
        const file = fs.createWriteStream(filePath);
        const fileName = path.basename(filePath);

        const cleanupAndReject = (err: Error) => {
            try { file.close(); } catch (_) { }
            fs.unlink(filePath, () => {
                flushBuffer();
                reject(err);
            });
        };

        const request = (currentUrl: string, redirectsLeft: number) => {
            const req = https.request(currentUrl, (response) => {
                const status = response.statusCode ?? 0;

                // handle redirects
                if (status >= 300 && status < 400 && response.headers.location) {
                    if (redirectsLeft === 0) {
                        response.destroy();
                        return cleanupAndReject(new Error('Too many redirects'));
                    }
                    const next = new URL(response.headers.location, currentUrl).toString();
                    response.destroy();
                    return request(next, redirectsLeft - 1);
                }

                if (status !== 200) {
                    response.destroy();
                    return cleanupAndReject(new Error(`Download failed, status ${status}`));
                }

                const total = parseInt((response.headers['content-length'] as string) || '0', 10);
                let downloaded = 0;

                response.on('data', (chunk: Buffer) => {
                    downloaded += chunk.length;
                    if (total) {
                        // TODO: extract progress bar logic
                        const percent = ((downloaded / total) * 100);
                        process.stdout.clearLine(0);
                        process.stdout.cursorTo(0);
                        process.stdout.write(`Downloading ${fileName}: [${'='.repeat(percent / 4)}${' '.repeat(25 - percent / 4)}] ${percent.toFixed(2)}%`);
                    } else {
                        process.stdout.clearLine(0);
                        process.stdout.cursorTo(0);
                        process.stdout.write(`Downloading ${fileName}: ${formatBytes(downloaded)} bytes`);
                    }
                });

                response.pipe(file);

                file.on('finish', () => {
                    file.close(() => {
                        flushBuffer();
                        resolve(filePath);
                    });
                });

                response.on('error', (err) => cleanupAndReject(err));
            });

            req.on('error', (err) => cleanupAndReject(err));
            req.end();
        };

        file.on('error', (err) => cleanupAndReject(err));

        request(url, maxRedirects);
    });
};

/**
 * extracts and archive using `tar -xf archivePath -C outDir`
 * @param archivePath
 * @param outDir
 */
export function extract(archivePath: string, outDir: string): Promise<void> {
    return new Promise((resolve, reject) => {
        if (!fs.existsSync(outDir)) {
            // ensure dir exists
            fs.mkdirSync(outDir, { recursive: true });
        }

        execFile('tar', ['-xf', archivePath, '-C', outDir], (err, _, stderr) => {
            if (err) {
                return reject(new Error(`Extraction failed: ${stderr || err.message}`));
            }
            resolve();
        });
    });
}

/**
 * creates a symbolic link between `outDir/filename` and `~/.eddy.sh/bin/filename`
 * @param dir
 * @param filename
 */
export function symlink(dir: string, filename: string) {
    if (!fs.existsSync(EDDY_BIN_DIR)) {
        fs.mkdirSync(EDDY_BIN_DIR, { mode: 0o755, recursive: true });
    }

    const target = path.join(EDDY_BIN_DIR, filename);
    try {
        if (fs.existsSync(target)) {
            fs.unlinkSync(target);
        }
        fs.symlinkSync(path.join(dir, filename), target);
    } catch (err) {
        // TODO: handle error
        logger.error(`Failed to create symlink: ${err}`);
    }
}

export const chmod755 = (targetPath: string, filename: string) => {
    const bin = path.join(targetPath, filename);

    fs.chmod(bin, 0o755, (err) => {
        if (err) {
            // TODO: handle error
            logger.error(err);
        }
    });
};

/**
 * renames a directory like mv pathname/oldName pathname/newName;
 *
 * In the case that newPath already exists, it will be overwritten.
 * If there is a directory at newPath, an error will be raised instead.
 * @param pathname absolute path
 * @param oldName current dir name
 * @param newName new dir name
 */
export const rename = (pathname: string, oldName: string, newName: string) => {
    return new Promise<void>((res, rej) => {
        const newPath = path.join(pathname, newName);
        if (fs.existsSync(newPath)) {
            logger.warn(`${newPath} directory is not empty. if you wanted to reinstall same version, remove it first`);
            res();
        }
        fs.rename(path.join(pathname, oldName), newPath, (err) => {
            // TODO: handle error
            if (err) rej(err);
            res();
        });
    });
};

/**
 * performs rm -rf on given `path`;
 * await for `fs.rm` to complete
 */
export const remove = (path: string) => {
    return new Promise<void>((res, rej) => {
        fs.rm(path, { recursive: true, force: true }, (err) => {
            if (err) {
                // TODO: handle err
                rej(err);
            }
            res();
        });
    });
};

/**
 * Sometimes tools use version in their package name,
 * so in order to determine proper version, following redirect is needed
 *
 * think of a way to test it reliably
 * publish eddy.sh version and use version from package.json maybe?
 * @param url to follow the redirect
 */
export async function resolveLatestVersion(url: string): Promise<semver> {
    return new Promise((resolve, reject) => {
        const req = https.request(url, { method: "HEAD" }, (res) => {
            const location = res.headers.location;
            if (location) {
                const match = location.match(/(\d+\.\d+\.\d+)/);
                if (match?.[1]) {
                    resolve(match[1] as semver);
                } else {
                    reject(new Error("Could not extract version and filename from redirect URL"));
                }
            } else {
                reject(new Error("No redirect location header found"));
            }
        });
        req.on("error", reject);
        req.end();
    });
}


export function formatBytes(bytes: number): string {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

function flushBuffer() { console.log('\n'); };