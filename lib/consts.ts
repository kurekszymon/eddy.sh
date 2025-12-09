import os from 'os';
import path from 'path';

const homeDir = os.homedir();

const EDDY_DIR = path.join(homeDir, '.eddy.sh');
const EDDY_BIN_DIR = path.join(EDDY_DIR, 'bin');

export { EDDY_DIR, EDDY_BIN_DIR };
