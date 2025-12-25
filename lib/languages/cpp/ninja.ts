import type { Tool } from "@/lib/types";

/**
 * ninja tool shape; call like `ninja('1.13.2')`
 * @param version version of a lib, pass semantic version without leading `v`
 *
 * @link https://github.com/ninja-build/ninja/releases/download/v1.13.2/ninja-mac.zip
 */
export const ninja = (version: Tool['version']): Tool => ({
    name: 'ninja',
    version,
    lang: 'cpp',

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
});
