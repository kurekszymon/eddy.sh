import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";
import { ToolBlueprint } from "@/lib/languages/blueprint";

describe('cpp/conan', async () => {
    const cpp = await import("@/lib/languages/cpp/conan");
    const conan = cpp.conan('2.23.0');
    const blueprint = new ToolBlueprint(conan);

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(conan.pkgName).toBe('conan-2.23.0-macos-arm64.tgz');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(conan.pkgName).toBe('conan-2.23.0-windows-x86_64.zip');
    });

    test('picks latest version', async () => {
        const cpp = await import("@/lib/languages/cpp/conan");
        const conan = cpp.conan('latest');
        expect(conan.url).toBe(`https://github.com/conan-io/conan/releases/latest/download/${conan.pkgName}`);
    });

    test("downloads conan", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        expect(conan.url).toBe(`https://github.com/conan-io/conan/releases/download/2.23.0/${conan.pkgName}`);

        const dir = ensureToolDir(`cpp/conan/${conan.version}`, { check: true });
        await blueprint.download();

        expect(fs.existsSync(path.join(dir, conan.pkgName))).toBe(true);
    });

    test("installs conan", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir(`cpp/conan/${conan.version}`, { check: true });

        await blueprint.install();
        blueprint.use();

        const symlinkPath = path.join(EDDY_BIN_DIR, process.platform === 'win32' ? conan.pkgName : 'conan');
        const symlinkStats = fs.lstatSync(symlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        expect(target).toBe(path.join(dir, conan.customBinPath!, conan.name));
    });

    test("deletes conan installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        const dir = ensureToolDir(`cpp/conan/${conan.version}`, { check: true });
        await blueprint.install();
        expect(fs.existsSync(path.join(dir, conan.customBinPath!, conan.name))).toBe(true);
        expect(fs.existsSync(path.join(dir, conan.pkgName))).toBe(true);

        await blueprint.delete();
        expect(fs.existsSync(path.join(dir, conan.customBinPath!, conan.name))).toBe(false);
        expect(fs.existsSync(path.join(dir, conan.pkgName))).toBe(false);
    });
});