import fs from 'fs';
import path from 'path';

import type { Tool } from "@/lib/types";
import {
    ensureToolDir,
    downloadFile,
    extract,
    symlink,
    chmod755,
    resolveLatestVersion
} from '@/lib/shared';

/**
 * ninja tool shape; call like `ninja('1.13.2')`
 * @param version version of a lib, pass semantic version without leading `v`
 *
 * @link https://github.com/ninja-build/ninja/releases/download/v1.13.2/ninja-mac.zip
 */
export const ninja = (version: Tool['version']): Tool => ({
    name: 'ninja',
    version,

    get pkgName() {
        if (process.platform === 'win32') {
            return 'ninja-win.zip';
        }
        if (process.platform === 'darwin') {
            return 'ninja-mac.zip';
        }

        throw new Error("Platform not supported!");
    },
    get url() {
        if (version === 'latest') {
            return `https://github.com/ninja-build/ninja/releases/latest/download/${this.pkgName}`;
        }
        return `https://github.com/ninja-build/ninja/releases/download/v${this.version}/${this.pkgName}`;
    },

    async install() {
        const outDir = ensureToolDir(`cpp/ninja/${this.version}`);

        if (version === 'latest') {
            this.version = await resolveLatestVersion(this.url);
        }

        const archivePath = await this.download();
        await extract(archivePath, outDir);
    },
    async download() {
        const dir = ensureToolDir('cpp/ninja');
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, this.url);
        return filePath;
    },
    use() {
        const ninjaDir = ensureToolDir(`cpp/ninja/${this.version}`);

        if (!fs.existsSync(ninjaDir)) {
            throw new Error(`${this.name}@${this.version} is not installed yet.`);
        }

        chmod755(ninjaDir, this.name);
        symlink(ninjaDir, this.name);
    }
});
