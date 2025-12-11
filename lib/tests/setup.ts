import { mock, beforeEach, afterEach } from "bun:test";

import fs from 'fs';
import path from 'path';
import { tmpdir } from 'os';

const tmpDir = path.join(tmpdir(), "eddy-test");

mock.module("os", () => ({
    homedir: () => tmpDir,
}));

beforeEach(() => {
    fs.mkdtempSync(tmpDir);
});

afterEach(() => {
    fs.rmSync(tmpDir, { recursive: true, force: true });
});
