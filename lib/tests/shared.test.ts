import { expect, test } from "bun:test";

import fs from "fs";
import path from "path";

import { tmpDir } from "./setup";


test("creates files in temp home", async () => {
    const { createToolDir } = await import("../shared");

    const dirName = 'test';
    createToolDir(dirName);

    const dir = path.join(tmpDir, '.eddy.sh', dirName);

    expect(fs.existsSync(dir)).toBe(true);
});
