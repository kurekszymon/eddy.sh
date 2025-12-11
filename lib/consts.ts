// TODO: investigate - needed so in tests value is correctly mocked
// if you use import os from 'os' it will not be mocked correctly in tests.
import { homedir } from 'os';
import path from 'path';

const homeDir = homedir();

const EDDY_DIR = path.join(homeDir, '.eddy.sh');
const EDDY_BIN_DIR = path.join(EDDY_DIR, 'bin');

export { EDDY_DIR, EDDY_BIN_DIR };
