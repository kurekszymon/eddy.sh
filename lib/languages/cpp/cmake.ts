import path from 'path';
import fs from 'fs';

import {
    downloadFile,
    ensureToolDir,
    extract,
    remove,
    rename,
    resolveLatestVersion,
    symlink
} from "@/lib/shared";
import type { Tool } from "@/lib/types";


export const getBasePkgName = (pkgName: string) => path.basename(pkgName).replace(/\.(tar\.gz|zip)$/, '');

export const CMAKE_BIN_PATH = process.platform === 'darwin'
    ? 'CMake.app/Contents/bin'
    : 'bin';
/**
 * cmake tool shape; call like `cmake('4.1.4')` or `cmake('latest')`
 * @param version version of a lib, pass semantic version without leading `v`
 *
 * @link https://github.com/Kitware/CMake/releases/download/v4.1.4/cmake-4.1.4-macos-universal.tar.gz
 */
export const cmake = (version: Tool['version']): Tool => ({
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
    get url() {
        if (version === 'latest') {
            return `https://github.com/Kitware/CMake/releases/latest/download/${this.pkgName}`;
        }

        return `https://github.com/Kitware/CMake/releases/download/v${this.version}/${this.pkgName}`;
    },

    async download() {
        const dir = ensureToolDir('cpp/cmake');
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, this.url);
        return filePath;
    },
    async install() {
        const cmakeDir = ensureToolDir('cpp/cmake');

        if (version === 'latest') {
            this.version = await resolveLatestVersion(this.url);
        }

        const archivePath = await this.download();
        await extract(archivePath, cmakeDir);

        await rename(cmakeDir, getBasePkgName(this.pkgName), this.version);
    },
    use() {
        const cmakeDir = ensureToolDir('cpp/cmake');

        const binDir = path.join(cmakeDir, this.version, CMAKE_BIN_PATH);

        if (!fs.existsSync(binDir)) {
            throw new Error(`${this.name}@${this.version} is not installed yet.`);
        }

        ['ccmake', 'cmake', 'cpack', 'ctest'].forEach(bin => {
            symlink(binDir, bin);
        });
    },
    async delete() {
        const cmakeDir = ensureToolDir(`cpp/cmake/${this.version}`);
        await remove(cmakeDir);
    }
});