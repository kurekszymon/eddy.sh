export type semver = `${number}.${number}.${number}`;

export interface Tool {
    name: string;
    version: 'latest' | semver;

    get pkgName(): string;
    get url(): string;

    download: () => Promise<string>;
    install: () => Promise<void>;
    delete: () => Promise<void>;
    use: () => void;
}