import { type Tool as ITool, type IToolBlueprint } from '@/lib/types';
import { downloadFile, ensureToolDir, extract, getBasePkgName, remove, rename, resolveLatestVersion, symlink } from '../shared';

import path from 'path';
import fs from 'fs';
import { logger } from '../logger';

// TODO: improve logging

export class ToolBlueprint implements IToolBlueprint {
    private url: string;
    private name: string;
    private pkgName: string;
    private lang: ITool['lang'];
    private version: ITool['version'];

    private links: string[] = [];
    private customBinPath: string = '';
    private renameNested: boolean = false;

    constructor(
        name: string,
        pkgName: string,
        url: string,
        lang: ITool['lang'],
        version: ITool['version'],
        links?: string[],
        customBinPath?: string,
        renameNested?: boolean,
    ) {
        this.url = url;
        this.name = name;
        this.pkgName = pkgName;
        this.version = version;
        this.lang = lang;

        this.links = links || [];
        this.customBinPath = customBinPath || '';
        this.renameNested = Boolean(renameNested);
    }

    async download(): Promise<string> {
        const dir = ensureToolDir(`${this.lang}/${this.name}`);
        const filePath = path.join(dir, this.pkgName);

        await downloadFile(filePath, this.url);
        return filePath;
    };

    async install(): Promise<void> {
        if (this.version === 'latest') {
            this.version = await resolveLatestVersion(this.url);
        }

        const dir = ensureToolDir(`${this.lang}/${this.name}/${this.version}`);

        const archivePath = await this.download();

        await extract(archivePath, dir);
        if (this.renameNested) {
            await rename(dir, getBasePkgName(this.pkgName), this.version);
        }
    };

    use(): void {
        const dir = ensureToolDir(`${this.lang}/${this.name}`, { check: true });
        const binDir = path.join(dir, this.version, this.customBinPath || '');

        if (!fs.existsSync(binDir)) {
            throw new Error(`${this.name}@${this.version} is not installed yet.`);
        }

        if (this.links.length > 0) {
            return this.links.forEach(bin => symlink(binDir, bin));
        }

        return symlink(binDir, this.name);
    };

    async delete(): Promise<void> {
        const dir = ensureToolDir(`${this.lang}/${this.name}/${this.version}`, { check: true });
        const archive = ensureToolDir(`${this.lang}/${this.name}/${this.pkgName}`, { check: true });

        const result = await Promise.allSettled([remove(archive), remove(dir)]);

        // TODO: notify about failed processes
        if (result.some(r => r.status === 'rejected')) {
            logger.info(`Failed to delete ${this.name}@${this.version}`);
        } else {
            logger.info(`Successfully deleted ${this.name}@${this.version}`);
        }
    };
}