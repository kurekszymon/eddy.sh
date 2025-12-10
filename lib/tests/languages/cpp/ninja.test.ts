import { describe, expect, test } from "bun:test";

import fs from 'fs';

describe('cpp/ninja', async () => {
    const cpp = await import("@/lib/languages/cpp/ninja");
    const ninja = cpp.ninja('1.13.2');

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(ninja.pkgName).toBe('ninja-mac.zip');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(ninja.pkgName).toBe('ninja-win.zip');
    });

    test("installs ninja", async () => {
        const dir = await ninja.download();

        expect(fs.existsSync(dir)).toBe(true);
    });
});