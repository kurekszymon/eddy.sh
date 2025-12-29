import type { IToolInfo, ToolVersion } from '@/lib/types';
import { getLatestGoVersion } from './fetchLatest';

const latest = await getLatestGoVersion();

export const lang = (version: ToolVersion): IToolInfo => ({
    name: 'go-language',
    lang: 'go',
    version,

    steps: ['extract'],
    links: ['go', 'gofmt'],
    customBinPath: 'go/bin',

    get pkgName() {
        const platform = process.platform;

        if (platform === 'darwin') {
            return `go${this.version}.darwin-arm64.tar.gz`;
        }
        if (platform === 'win32') {
            return `go${this.version}.windows-386.zip`;
        }

        throw new Error('Unsupported platform/arch');
    },

    get url() {
        if (this.version === 'latest') {
            this.version = latest;
        }

        return `https://go.dev/dl/${this.pkgName}`;
    },
});