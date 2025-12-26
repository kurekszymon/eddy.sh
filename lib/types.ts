export type semver = `${number}.${number}.${number}`;

export type InstallStep = 'extract' | 'rename' | 'chmod';

export interface Tool {
    name: string;
    version: 'latest' | semver;

    get pkgName(): string;
    get url(): string;

    lang: 'cpp' | 'go',
    customBinPath?: string;
    links?: string[];
    steps: InstallStep[];

    download?: () => Promise<string>;
    install?: () => Promise<void>;
    delete?: () => Promise<void>;
    use?: () => void;
}

export interface IToolBlueprint {
    download: () => Promise<string>;
    install: (opts: { renameNested: boolean; }) => Promise<void>;
    delete: () => Promise<void>;
    use: () => void;
}