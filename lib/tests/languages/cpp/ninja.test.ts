import { EDDY_BIN_DIR } from "@/lib/consts";
import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

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
        const { ensureToolDir } = await import("@/lib/shared");

        const dir = ensureToolDir('cpp/ninja');
        await ninja.install();

        const symlinkPath = path.join(EDDY_BIN_DIR, 'ninja');
        const symlinkStats = fs.lstatSync(symlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        expect(target).toBe(path.join(dir, 'ninja'));
    });
});