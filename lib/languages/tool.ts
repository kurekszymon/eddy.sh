import { type Tool as ITool, type IToolBlueprint } from '@/lib/types';
import { downloadFile, ensureToolDir, extract, remove, resolveLatestVersion, symlink } from '../shared';

import path from 'path';
import fs from 'fs';
import { logger } from '../logger';

// TODO: improve logging

export class ToolBlueprint implements IToolBlueprint {
    private info: ITool;

    constructor(info: ITool) {
        this.info = info;
    }

    async download(): Promise<string> {
        const { lang, name, pkgName, url } = this.info;

        const dir = ensureToolDir(`${lang}/${name}`);
        const filePath = path.join(dir, pkgName);

        await downloadFile(filePath, url);
        return filePath;
    };

    async install(): Promise<void> {
        const { lang, name, version, url } = this.info;

        if (version === 'latest') {
            this.info.version = await resolveLatestVersion(url);
        }

        const dir = ensureToolDir(`${lang}/${name}/${version}`);

        const archivePath = await this.download();

        await extract(archivePath, dir);
    };

    use(): void {
        const { lang, name, version, links, customBinPath } = this.info;

        const dir = ensureToolDir(`${lang}/${name}`, { check: true });
        const binDir = path.join(dir, version, customBinPath || '');

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