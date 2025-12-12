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

    test('picks latest version', async () => {
        const cpp = await import("@/lib/languages/cpp/cmake");
        const cmake = cpp.cmake('latest');

        expect(cmake.url).toBe(`https://github.com/Kitware/CMake/releases/latest/download/${cmake.pkgName}`);
    });

    test("installs cmake", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        expect(cmake.url).toBe(`https://github.com/Kitware/CMake/releases/download/v4.1.4/${cmake.pkgName}`);

        const dir = ensureToolDir('cpp/cmake');
        await cmake.download();

        expect(fs.existsSync(path.join(dir, cmake.pkgName))).toBe(true);
    });
});