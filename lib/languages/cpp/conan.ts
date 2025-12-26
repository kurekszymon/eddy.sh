import type { IToolInfo, ToolVersion } from "@/lib/types";

export const conan = (version: ToolVersion): IToolInfo => ({
    name: 'conan',
    lang: 'cpp',
    version,

    steps: ['extract'],
    customBinPath: 'bin',

    get pkgName() {
        if (process.platform === 'win32') {
            return `conan-${this.version}-windows-x86_64.zip`;
        }
        if (process.platform === 'darwin') {
            return `conan-${this.version}-macos-arm64.tgz`;
        }

        throw new Error("Platform not supported!");
    },
    get url() {
        if (version === 'latest') {
            return `https://github.com/conan-io/conan/releases/latest/download/${this.pkgName}`;
        }
        return `https://github.com/conan-io/conan/releases/download/${this.version}/${this.pkgName}`;
    },
});