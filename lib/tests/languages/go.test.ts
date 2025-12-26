import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";
import { ToolBlueprint } from "@/lib/languages/blueprint";

describe('go', async () => {
    const goMod = await import("@/lib/languages/go/lang");
    const go = goMod.lang('1.25.5');
    const blueprint = new ToolBlueprint(go);

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(go.pkgName).toBe('go1.25.5.darwin-arm64.tar.gz');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(go.pkgName).toBe('go1.25.5..windows-386.zip');
    });

    test.skip('picks latest version', async () => {
    });

    test("downloads go", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        expect(go.url).toBe(`https://go.dev/dl/${go.pkgName}`);

        const dir = ensureToolDir(`go/go-language/${go.version}`, { check: true });
        await blueprint.download();

        expect(fs.existsSync(path.join(dir, go.pkgName))).toBe(true);
    });

    test("installs go", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir(`go/go-language/${go.version}`, { check: true });

        await blueprint.install();
        blueprint.use();

        const symlinkPath = path.join(EDDY_BIN_DIR, 'go');
        const fmtSymlinkPath = path.join(EDDY_BIN_DIR, 'gofmt');
        const symlinkStats = fs.lstatSync(symlinkPath);
        const fmtSymlinkStats = fs.lstatSync(fmtSymlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);
        expect(fmtSymlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        const targetFmt = fs.readlinkSync(fmtSymlinkPath);

        expect(target).toBe(path.join(dir, go.customBinPath!, 'go'));
        expect(targetFmt).toBe(path.join(dir, go.customBinPath!, 'gofmt'));
    });

    test("deletes go installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        const dir = ensureToolDir(`go/go-language/${go.version}`, { check: true });

        await blueprint.install();
        expect(fs.existsSync(path.join(dir, go.customBinPath!))).toBe(true);
        expect(fs.existsSync(path.join(dir, go.pkgName))).toBe(true);

        await blueprint.delete();
        expect(fs.existsSync(path.join(dir, go.customBinPath!))).toBe(false);
        expect(fs.existsSync(path.join(dir, go.pkgName))).toBe(false);
    });
});