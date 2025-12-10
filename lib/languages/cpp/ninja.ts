import path from 'path';

import type { Tool } from "@/lib/types";
import { createToolDir, downloadFile } from '@/lib/shared';

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
        const filePath = await this.download();
        // handle install;
    },
    download: async function () {
        const dir = createToolDir('cpp/ninja');
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, `https://github.com/ninja-build/ninja/releases/download/v${this.version}/${this.pkgName}`);
        return filePath;
    },

});
