import type { IToolInfo, ToolVersion } from "@/lib/types";

export const bazel = (version: ToolVersion): IToolInfo => ({
    name: 'bazel',
    lang: 'cpp',
    version,

    steps: ['rename', 'chmod'],

    get pkgName() {
        if (process.platform === 'win32') {
            return `bazel-${this.version}-windows-x86_64.exe`;
        }
        if (process.platform === 'darwin') {
            return `bazel-${this.version}-darwin-arm64`;
        }
        throw new Error("Platform not supported!");
    },
    get url() {
        if (version === 'latest') {
            return `https://github.com/bazelbuild/bazel/releases/latest/download/${this.pkgName}`;
        }
        return `https://github.com/bazelbuild/bazel/releases/download/${this.version}/${this.pkgName}`;
    },
});