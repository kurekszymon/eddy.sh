export interface Tool {
    name: string;
    version: 'latest' | `${number}.${number}.${number}`;

    get pkgName(): string;
    get url(): string;

    install: () => Promise<void>;
    download: () => Promise<string>;
}