import path from 'path';
import fs from 'fs';

import type { Tool } from "@/lib/types";
import {
    ensureToolDir,
    downloadFile,
    chmod755,
    symlink,
    resolveLatestVersion,
    remove,
    extract,
} from '@/lib/shared';

// Conan is distributed as a single Python script or as a standalone binary for some platforms.
// We'll use the official standalone installer for Linux/macOS, and the .exe for Windows.

export const conan = (version: Tool['version']): Tool => ({
    name: 'conan',
    version,

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

    async download() {
        const dir = ensureToolDir('cpp/conan');
        const filePath = path.join(dir, this.pkgName);
        await downloadFile(filePath, this.url);
        return filePath;
    },
    async install() {
        if (version === 'latest') {
            this.version = await resolveLatestVersion(this.url);
        }

        const conanDir = ensureToolDir(`cpp/conan/${this.version}`);

        const archivePath = await this.download();
        await extract(archivePath, conanDir);
    },
    use() {
        const conanDir = ensureToolDir(`cpp/conan/${this.version}`, { check: false });

        if (!fs.existsSync(conanDir)) {
            // remove double exists sync version, replace with only one
            throw new Error(`${this.name}@${this.version} is not installed yet.`);
        }

        chmod755(conanDir, this.pkgName);
        symlink(path.join(conanDir, 'bin'), this.name);
    },
    async delete() {
        const conanDir = ensureToolDir(`cpp/conan/${this.version}`);
        await remove(conanDir);
    }
});