import { describe, expect, test } from "bun:test";

import fs, { readdirSync } from 'fs';
import path from 'path';

describe('cpp/cmake', async () => {
    const cpp = await import("@/lib/languages/cpp/cmake");
    const cmake = cpp.cmake('4.1.4');

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(cmake.pkgName).toBe('cmake-4.1.4-macos-universal.tar.gz');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(cmake.pkgName).toBe('cmake-4.1.4-windows-x86_64.zip');
    });

    test("installs cmake", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        const dir = ensureToolDir('cpp/cmake');
        await cmake.download();

        expect(fs.existsSync(path.join(dir, cmake.pkgName))).toBe(true);
    });
});