import { type InstallStep, type IToolInfo, type IToolBlueprint } from '@/lib/types';
import {
    downloadFile,
    ensureToolDir,
    extract,
    rename,
    remove,
    resolveLatestVersion,
    symlink,
    chmod755
} from '@/lib/shared';
import { logger } from '@/lib/logger';

import path from 'path';
import fs from 'fs';

// TODO: improve logging

export class ToolBlueprint implements IToolBlueprint {
    private info: IToolInfo;

    constructor(info: IToolInfo) {
        this.info = info;
    }

    async download(): Promise<string> {
        const { lang, name, pkgName, version, url } = this.info;

        const dir = ensureToolDir(`${lang}/${name}/${version}`);
        const filePath = path.join(dir, pkgName);

        await downloadFile(filePath, url);
        return filePath;
    };

    async install(): Promise<void> {
        const { name, version, url, pkgName, steps } = this.info;

        if (version === 'latest') {
            this.info.version = await resolveLatestVersion(url);
        }

        const archivePath = await this.download();
        const dir = path.dirname(archivePath);

        const stepHandlers: Record<string, () => Promise<void>> = {
            extract: async () => { await extract(archivePath, dir); },
            rename: async () => { await rename(dir, pkgName, name); },
            chmod: async () => { chmod755(dir, name); },
        };

        for (const step of steps) {
            const handler = stepHandlers[step];
            if (handler) {
                await handler();
            }
        }
    };

    use(): void {
        const { lang, name, version, links, customBinPath } = this.info;

        const dir = ensureToolDir(`${lang}/${name}/${version}`, { check: true });

        const binDir = path.join(dir, customBinPath || '');

        if (!fs.existsSync(binDir)) {
            throw new Error(`${name}@${version} is not installed yet.`);
        }

        if (links && links.length > 0) {
            return links.forEach(bin => symlink(binDir, bin));
        }

        return symlink(binDir, name);
    };

    async delete(): Promise<void> {
        const { lang, name, version, pkgName } = this.info;

        const dir = ensureToolDir(`${lang}/${name}/${version}`, { check: true });
        const archive = ensureToolDir(`${lang}/${name}/${pkgName}`, { check: true });

        const result = await Promise.allSettled([remove(archive), remove(dir)]);

        // TODO: notify about failed processes
        if (result.some(r => r.status === 'rejected')) {
            logger.info(`Failed to delete ${name}@${version}`);
        } else {
            logger.info(`Successfully deleted ${name}@${version}`);
        }
    };
}