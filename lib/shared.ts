import fs from 'fs';
import path from 'path';
import https from 'https';

import { EDDY_DIR } from '@/lib/consts';
import { logger } from '@/lib/logger';

export const createToolDir = (dirName: string) => {
    const dir = path.join(EDDY_DIR, dirName);

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
        const fileName = filePath.split('/').at(-1);

        const cleanupAndReject = (err: Error) => {
            try { file.close(); } catch (_) { }
            fs.unlink(filePath, () => {
                console.log('\n');
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
                        if (process.stdout?.isTTY) {
                            process.stdout.clearLine(0);
                            process.stdout.cursorTo(0);
                            process.stdout.write(`Downloading ${fileName}: [${'='.repeat(percent / 4)}${' '.repeat(25 - percent / 4)}] ${percent}%`);
                        } else {
                            logger.info(`Downloading ${fileName}: ${percent.toFixed(2)}%`);
                        }
                    } else {
                        if (process.stdout?.isTTY) {
                            process.stdout.clearLine(0);
                            process.stdout.cursorTo(0);
                            process.stdout.write(`Downloading ${fileName}: ${formatBytes(downloaded)} bytes`);
                        } else {
                            logger.info(`Downloading ${fileName}: ${formatBytes(downloaded)} bytes`);
                        }
                    }
                });

                response.pipe(file);

                file.on('finish', () => {
                    file.close(() => {
                        resolve(filePath);
                        console.log('\n');
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

export function formatBytes(bytes: number): string {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}