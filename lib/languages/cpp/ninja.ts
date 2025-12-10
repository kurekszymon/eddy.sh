import path from 'path';
import fs from 'fs';

import type { Tool } from "@/lib/types";
import { ensureToolDir, downloadFile, extract, symlink } from '@/lib/shared';
import { logger } from '@/lib/logger';

/**
 * ninja tool shape; call like `ninja('1.13.2')`
 * @param version version of a lib, pass semantic version without leading `v`
 *
 * @link https://github.com/ninja-build/ninja/releases/download/v1.13.2/ninja-mac.zip
 */
export const ninja = (version: string): Tool => ({
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
    install: async function () {
        const outDir = ensureToolDir('cpp/ninja'); // TODO: infer lang/name
        const archivePath = await this.download();
        await extract(archivePath, outDir);

        const bin = path.join(outDir, this.name);
        fs.chmod(bin, 0o755, (err) => {
            if (err) {
                // TODO: handle error
                logger.error(err);
            }
        });

        symlink(outDir, this.name);
    },
    download: async function () {
        const dir = ensureToolDir('cpp/ninja');
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, `https://github.com/ninja-build/ninja/releases/download/v${this.version}/${this.pkgName}`);
        return filePath;
    },
});
