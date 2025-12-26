export type semver = `${number}.${number}.${number}`;

export type InstallStep = 'extract' | 'rename' | 'chmod';

export type ToolVersion = 'latest' | semver;

export interface IToolInfo {
    name: string;
    version: ToolVersion;

    get pkgName(): string;
    get url(): string;

    lang: 'cpp' | 'go',
    customBinPath?: string;
    links?: string[];
    steps: InstallStep[];
}

export interface IToolBlueprint {
    download: () => Promise<string>;
    install: () => Promise<void>;
    delete: () => Promise<void>;
    use: () => void;
}