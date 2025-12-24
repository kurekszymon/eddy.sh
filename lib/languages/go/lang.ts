import path from 'path';
import type { Tool } from '@/lib/types';
import {
    ensureToolDir,
    downloadFile,
    extract,
    symlink,
    remove,
} from '@/lib/shared';
import { logger } from '@/lib/logger';

export const go = (version: Tool['version']): Tool => ({
    name: 'go',
    version,

    get pkgName() {
        const platform = process.platform;

        if (platform === 'darwin') {
            return `go${this.version}.darwin-arm64.tar.gz`;
        }
        if (platform === 'win32') {
            return `go${this.version}.windows-386.zip`;
        }

        throw new Error('Unsupported platform/arch');
    },

    get url() {
        if (this.version === 'latest') {
            // Go does not have a 'latest' download URL, so resolve it first
            throw new Error('Please resolve the latest Go version before downloading');
        }
        return `https://go.dev/dl/${this.pkgName}`;
    },

    async download() {
        const dir = ensureToolDir('go');
        const filePath = path.join(dir, this.pkgName);
        await downloadFile(filePath, this.url);
        return filePath;
    },

    async install() {
        const goDir = ensureToolDir(`go/${this.version}`);
        const archivePath = await this.download();
        await extract(archivePath, goDir);
    },

    use() {
        const goDir = ensureToolDir(`go/${this.version}/go`, { check: true });

        symlink(path.join(goDir, 'bin'), 'go');
        symlink(path.join(goDir, 'bin'), 'gofmt');
    },

    async delete() {
        const goDir = ensureToolDir(`go/${this.version}`, { check: true });
        const goArchive = ensureToolDir(`go/${this.pkgName}`, { check: true });

        logger.info(`Deleting ${this.name}@${this.version}`);

        const result = await Promise.allSettled([remove(goArchive), remove(goDir)]);

        if (result.some(r => r.status === 'rejected')) {
            logger.info(`Failed to delete ${this.name}@${this.version}`);
        } else {
            logger.info(`Successfully deleted ${this.name}@${this.version}`);
        }
    }
});