import { describe, expect, test } from "bun:test";

import fs from 'fs';
import path from 'path';

import { EDDY_BIN_DIR } from "@/lib/consts";

describe('cpp/bazel', async () => {
    const cpp = await import("@/lib/languages/cpp/bazel");
    const bazel = cpp.bazel('8.5.0');

    test.if(process.platform === 'darwin')("checks pkgName", () => {
        expect(bazel.pkgName).toBe('bazel-8.5.0-darwin-arm64');
    });

    test.if(process.platform === 'win32')("checks pkgName", () => {
        expect(bazel.pkgName).toBe('bazel-8.5.0-windows-x86_64.exe');
    });

    test.if(process.platform === 'linux')("checks pkgName", () => {
        expect(bazel.pkgName).toBe('bazel-8.5.0-linux-x86_64');
    });

    test('picks latest version', async () => {
        const cpp = await import("@/lib/languages/cpp/bazel");
        const bazel = cpp.bazel('latest');
        expect(bazel.url).toBe(`https://github.com/bazelbuild/bazel/releases/latest/download/${bazel.pkgName}`);
    });

    test("downloads bazel", async () => {
        const { ensureToolDir } = await import("@/lib/shared");

        expect(bazel.url).toBe(`https://github.com/bazelbuild/bazel/releases/download/8.5.0/${bazel.pkgName}`);
        const dir = ensureToolDir('cpp/bazel');

        await bazel.download();
        expect(fs.existsSync(path.join(dir, bazel.version))).toBe(true);
    });

    test("installs bazel", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir('cpp/bazel');
        await bazel.install();
        bazel.use();

        const symlinkPath = path.join(EDDY_BIN_DIR, process.platform === 'win32' ? bazel.pkgName : 'bazel');
        const symlinkStats = fs.lstatSync(symlinkPath);
        expect(symlinkStats.isSymbolicLink()).toBe(true);

        const target = fs.readlinkSync(symlinkPath);
        expect(target).toBe(path.join(dir, bazel.version, bazel.name));
    });

    test("deletes bazel installation", async () => {
        const { ensureToolDir } = await import("@/lib/shared");
        const dir = ensureToolDir('cpp/bazel');

        await bazel.install();
        expect(fs.existsSync(path.join(dir, bazel.version))).toBe(true);

        await bazel.delete();
        expect(fs.existsSync(path.join(dir, bazel.version))).toBe(false);
    });
});