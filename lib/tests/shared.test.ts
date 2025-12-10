import { expect, test } from "bun:test";

import fs from "fs";
import path from 'path';

test("creates files in temp home", async () => {
    const { ensureToolDir } = await import("@/lib/shared");

    const dir = ensureToolDir('test');

    expect(fs.existsSync(dir)).toBe(true);
});

test("downloads file", async () => {
    const { ensureToolDir, downloadFile } = await import("@/lib/shared");

    const dir = ensureToolDir('test');
    const filePath = path.join(dir, 'Makefile');

    await downloadFile(filePath, 'https://github.com/kurekszymon/eddy.sh/blob/main/Makefile');

    expect(fs.existsSync(filePath)).toBe(true);
});

test("extracts archive", async () => {
    const { extract, ensureToolDir } = await import("@/lib/shared");
    const dir = ensureToolDir('extract-test');

    await extract("./lib/tests/fixtures/archive.zip", dir);
    const extractedFile = path.join(dir, "archive.txt");
    expect(fs.existsSync(extractedFile)).toBe(true);

    const content = fs.readFileSync(extractedFile, "utf8");
    expect(content.trim()).toBe("hello extract");
});

test("creates symlink that points to correct file", async () => {
    const { symlink, ensureToolDir } = await import("@/lib/shared");
    const { EDDY_BIN_DIR } = await import("@/lib/consts");

    const dir = ensureToolDir('symlink-test');
    const filename = "dummy.txt";
    const filePath = path.join(dir, filename);
    fs.writeFileSync(filePath, "symlink test");

    symlink(dir, filename);

    const linkPath = path.join(EDDY_BIN_DIR, filename);
    const realPath = fs.readlinkSync(linkPath);
    expect(realPath).toBe(filePath);
});

test("formats bytes correctly", async () => {
    const { formatBytes } = await import("@/lib/shared");

    expect(formatBytes(512)).toBe("512 B");
    expect(formatBytes(1024)).toBe("1.0 KB");
    expect(formatBytes(1536)).toBe("1.5 KB");
    expect(formatBytes(1048576)).toBe("1.00 MB");
    expect(formatBytes(2097152)).toBe("2.00 MB");
});