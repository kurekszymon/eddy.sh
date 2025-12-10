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