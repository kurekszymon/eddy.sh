import path from 'path';

import { downloadFile, ensureToolDir } from "@/lib/shared";
import type { Tool } from "@/lib/types";

/**
 * cmake tool shape; call like `cmake('4.1.4')`
 * @param version version of a lib, pass semantic version without leading `v`
 *
 * @link https://github.com/Kitware/CMake/releases/download/v4.1.4/cmake-4.1.4-macos-universal.tar.gz
 */
export const cmake = (version: string): Tool => ({
    name: 'cmake',
    version,

    get pkgName() {
        if (process.platform === 'win32') {
            return `${this.name}-${this.version}-windows-x86_64.zip`;
        }
        if (process.platform === 'darwin') {
            return `${this.name}-${this.version}-macos-universal.tar.gz`;
        }

        throw new Error("Platform not supported!");
    },
    download: async function () {
        const dir = ensureToolDir('cpp/cmake'); // TODO infer from the file structure
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, `https://github.com/Kitware/CMake/releases/download/v${this.version}/${this.pkgName}`);
        return filePath;
    },
    async install() { }
});