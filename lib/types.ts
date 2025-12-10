export interface Tool {
    name: string;
    version: string;

    get pkgName(): string;

    install: () => Promise<void>;
    download: () => Promise<string>;
}