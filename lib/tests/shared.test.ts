import { expect, test } from "bun:test";

import fs from "fs";
import path from "path";

import { tmpDir } from "./setup";


test("creates files in temp home", async () => {
    const { createEddyDirs } = await import("../shared");
    createEddyDirs();

    const eddyDir = path.join(tmpDir, '.eddy.sh');
    expect(fs.existsSync(eddyDir)).toBe(true);

    const eddyBinDir = path.join(eddyDir, 'bin');
    expect(fs.existsSync(eddyBinDir)).toBe(true);
});
