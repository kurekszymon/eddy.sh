import { expect, test } from "bun:test";

import fs from "fs";
import path from 'path';

test("creates files in temp home", async () => {
    const { createToolDir } = await import("../shared");

    const dir = createToolDir('test');

    expect(fs.existsSync(dir)).toBe(true);
});

test("downloads file", async () => {
    const { createToolDir, downloadFile } = await import("../shared");

    const dir = createToolDir('test');
    const filePath = path.join(dir, 'Makefile');

    await downloadFile(filePath, 'https://github.com/kurekszymon/eddy.sh/blob/main/Makefile');

    expect(fs.existsSync(filePath)).toBe(true);
});

test("formats bytes correctly", async () => {
    const { formatBytes } = await import("../shared");

    expect(formatBytes(512)).toBe("512 B");
    expect(formatBytes(1024)).toBe("1.0 KB");
    expect(formatBytes(1536)).toBe("1.5 KB");
    expect(formatBytes(1048576)).toBe("1.00 MB");
    expect(formatBytes(2097152)).toBe("2.00 MB");
});