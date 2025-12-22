import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";

describe('go', async () => {
    const goMod = await import("@/lib/languages/go");
    const go = goMod.go('1.25.5');

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

        const dir = ensureToolDir('go', { check: true });
        await go.download();

        expect(fs.existsSync(path.join(dir, go.pkgName))).toBe(true);
    });

    test("installs go", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir('go', { check: true });

        await go.install();
        go.use();

        const symlinkPath = path.join(EDDY_BIN_DIR, 'go');
        const fmtSymlinkPath = path.join(EDDY_BIN_DIR, 'gofmt');
        const symlinkStats = fs.lstatSync(symlinkPath);
        const fmtSymlinkStats = fs.lstatSync(fmtSymlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);
        expect(fmtSymlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        const targetFmt = fs.readlinkSync(fmtSymlinkPath);
        expect(target).toBe(path.join(dir, go.version, go.name, 'bin', go.name));
        expect(targetFmt).toBe(path.join(dir, go.version, go.name, 'bin', 'gofmt'));
    });

    test("deletes go installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        const dir = ensureToolDir('go', { check: true });
        await go.install();
        expect(fs.existsSync(path.join(dir, go.version))).toBe(true);
        expect(fs.existsSync(path.join(dir, go.pkgName))).toBe(true);

        await go.delete();
        expect(fs.existsSync(path.join(dir, go.version))).toBe(false);
        expect(fs.existsSync(path.join(dir, go.pkgName))).toBe(false);
    });
});