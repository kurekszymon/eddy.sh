import type { Tool } from '@/lib/types';

export const lang = (version: Tool['version']): Tool => ({
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
            // Go does not have a 'latest' download URL, so resolve it first
            throw new Error('Please resolve the latest Go version before downloading');
        }
        return `https://go.dev/dl/${this.pkgName}`;
    },
});