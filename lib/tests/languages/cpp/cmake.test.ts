import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";
import { ToolBlueprint } from "@/lib/languages/blueprint";
import { getBasePkgName } from "@/lib/shared";

describe('cpp/cmake', async () => {
    const cpp = await import("@/lib/languages/cpp/cmake");
    const cmake = cpp.cmake('4.1.4');
    const blueprint = new ToolBlueprint(cmake);

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

    test("downloads cmake", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        expect(cmake.url).toBe(`https://github.com/Kitware/CMake/releases/download/v4.1.4/${cmake.pkgName}`);

        const dir = ensureToolDir(`cpp/cmake/${cmake.version}`, { check: true });
        await blueprint.download();

        expect(fs.existsSync(path.join(dir, cmake.pkgName))).toBe(true);
    });

    // TODO add test for latest version, so it always evaluates to `semantic.version.number` instead of latest
    test("installs cmake", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        expect(cmake.url).toBe(`https://github.com/Kitware/CMake/releases/download/v4.1.4/${cmake.pkgName}`);

        const dir = ensureToolDir(`cpp/cmake/${cmake.version}`, { check: true });
        await blueprint.install();
        blueprint.use();

        ['cmake', 'cpack', 'ctest', 'ccmake'].forEach(bin => {
            const symlinkPath = path.join(EDDY_BIN_DIR, bin);
            const symlinkStats = fs.lstatSync(symlinkPath);
            expect(symlinkStats.isSymbolicLink()).toBe(true);

            const target = fs.readlinkSync(symlinkPath);
            expect(target).toBe(path.join(dir, cmake.customBinPath!, bin));
        });
    }); // TODO: seperate tests between install and use

    test("deletes cmake installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir(`cpp/cmake/${cmake.version}`, { check: true });

        await blueprint.install();
        expect(fs.existsSync(path.join(dir, cmake.pkgName))).toBe(true);
        expect(fs.existsSync(path.join(dir, getBasePkgName(cmake.pkgName)))).toBe(true);

        await blueprint.delete();
        expect(fs.existsSync(path.join(dir, getBasePkgName(cmake.pkgName)))).toBe(false);
        expect(fs.existsSync(path.join(dir, cmake.pkgName))).toBe(false);
    });
});