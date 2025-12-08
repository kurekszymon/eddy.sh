import { mock, beforeAll, afterAll } from "bun:test";
import fs from 'fs';
import path from 'path';
import os from 'os';

export let tmpDir: string;

beforeAll(() => {
    tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "eddy-test-"));

    mock.module("os", () => ({
        homedir: () => tmpDir,
    }));
});

afterAll(() => {
    fs.rmSync(tmpDir, { recursive: true, force: true });
});
