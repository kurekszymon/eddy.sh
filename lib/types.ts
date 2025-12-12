export type semver = `${number}.${number}.${number}`;

export interface Tool {
    name: string;
    version: 'latest' | semver;

    get pkgName(): string;
    get url(): string;

    install: () => Promise<void>;
    download: () => Promise<string>;
}