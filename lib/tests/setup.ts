import { mock, beforeAll, afterAll } from "bun:test";

import fs from 'fs';
import path from 'path';
import os from 'os';

export const tmpDir = path.join(os.tmpdir(), "eddy-test");

mock.module("os", () => ({
    homedir: () => tmpDir,
}));

beforeAll(() => {
    fs.mkdtempSync(tmpDir);
});

afterAll(() => {
    fs.rmSync(tmpDir, { recursive: true, force: true });
});
