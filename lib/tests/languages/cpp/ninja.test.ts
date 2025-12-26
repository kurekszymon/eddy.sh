import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";
import { ToolBlueprint } from "@/lib/languages/blueprint";

describe('cpp/ninja', async () => {
    const cpp = await import("@/lib/languages/cpp/ninja");
    const ninja = cpp.ninja('1.13.2');
    const blueprint = new ToolBlueprint(ninja);

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(ninja.pkgName).toBe('ninja-mac.zip');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(ninja.pkgName).toBe('ninja-win.zip');
    });

    test('picks latest version', async () => {
        const cpp = await import("@/lib/languages/cpp/ninja");
        const ninja = cpp.ninja('latest');

        expect(ninja.url).toBe(`https://github.com/ninja-build/ninja/releases/latest/download/${ninja.pkgName}`);
    });

    test("downloads ninja", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        expect(ninja.url).toBe(`https://github.com/ninja-build/ninja/releases/download/v1.13.2/${ninja.pkgName}`);

        const dir = ensureToolDir(`cpp/ninja/${ninja.version}`, { check: true });
        await blueprint.download();

        expect(fs.existsSync(path.join(dir, ninja.pkgName))).toBe(true);
    });

    test("installs ninja", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        const dir = ensureToolDir(`cpp/ninja/${ninja.version}`, { check: true });
        await blueprint.install();
        blueprint.use();

        const symlinkPath = path.join(EDDY_BIN_DIR, ninja.name);
        const symlinkStats = fs.lstatSync(symlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        expect(target).toBe(path.join(dir, ninja.name));
    }); // TODO: seperate tests between install and use

    test("deletes ninja installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir(`cpp/ninja/${ninja.version}`, { check: true });

        await blueprint.install();
        expect(fs.existsSync(path.join(dir, ninja.name))).toBe(true);

        await blueprint.delete();
        expect(fs.existsSync(path.join(dir, ninja.name))).toBe(false);
    });
});