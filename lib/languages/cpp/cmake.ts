import type { Tool } from "@/lib/types";


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
    lang: 'cpp',
    links: ['ccmake', 'cmake', 'cpack', 'ctest'],
    customBinPath: CMAKE_BIN_PATH,

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
});