import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";
import { ToolBlueprint } from "@/lib/languages/blueprint";

describe('cpp/bazel', async () => {
    const cpp = await import("@/lib/languages/cpp/bazel");
    const bazel = cpp.bazel('8.5.0');
    const blueprint = new ToolBlueprint(bazel);

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(bazel.pkgName).toBe('bazel-8.5.0-darwin-arm64');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(bazel.pkgName).toBe('bazel-8.5.0-windows-x86_64.exe');
    });

    test('picks latest version', async () => {
        const cpp = await import("@/lib/languages/cpp/bazel");
        const bazel = cpp.bazel('latest');
        expect(bazel.url).toBe(`https://github.com/bazelbuild/bazel/releases/latest/download/${bazel.pkgName}`);
    });

    test("downloads bazel", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        expect(bazel.url).toBe(`https://github.com/bazelbuild/bazel/releases/download/8.5.0/${bazel.pkgName}`);
        const dir = ensureToolDir(`cpp/bazel/${bazel.version}`, { check: true });

        await blueprint.download();
        expect(fs.existsSync(path.join(dir, bazel.pkgName))).toBe(true);
    });

    test("installs bazel", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir(`cpp/bazel/${bazel.version}`, { check: true });

        await blueprint.install();
        blueprint.use();

        const symlinkPath = path.join(EDDY_BIN_DIR, process.platform === 'win32' ? bazel.pkgName : 'bazel');
        const symlinkStats = fs.lstatSync(symlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        expect(target).toBe(path.join(dir, bazel.name));
    });

    test("deletes bazel installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir(`cpp/bazel/${bazel.version}`, { check: true });

        await blueprint.install();
        expect(fs.existsSync(path.join(dir, bazel.name))).toBe(true);

        await blueprint.delete();
        expect(fs.existsSync(path.join(dir, bazel.name))).toBe(false);
    });
});