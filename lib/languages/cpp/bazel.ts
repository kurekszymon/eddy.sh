import path from 'path';
import fs from 'fs';

import type { Tool } from "@/lib/types";
import {
    ensureToolDir,
    downloadFile,
    symlink,
    chmod755,
    resolveLatestVersion,
    remove,
    rename
} from '@/lib/shared';

export const bazel = (version: Tool['version']): Tool => ({
    name: 'bazel',
    version,

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

    async download() {
        const dir = ensureToolDir(`cpp/bazel/${this.version}`);
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, this.url);
        return filePath;
    },
    async install() {
        if (version === 'latest') {
            this.version = await resolveLatestVersion(this.url);
        }

        const outDir = ensureToolDir(`cpp/bazel/${this.version}`);

        await this.download();
        await rename(outDir, this.pkgName, this.name);
    },
    use() {
        const bazelDir = ensureToolDir(`cpp/bazel/${this.version}`, { check: true });

        if (!fs.existsSync(bazelDir)) {
            // remove double exists sync version, replace with only one
            throw new Error(`${this.name}@${this.version} is not installed yet.`);
        }

        chmod755(bazelDir, this.name);
        symlink(bazelDir, this.name);
    },
    async delete() {
        const bazelDir = ensureToolDir(`cpp/bazel/${this.version}`, { check: true });
        await remove(bazelDir);
    }
});