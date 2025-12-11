// since the os module doesn't have a default export, it's only feasible to mock individual named exports
// attempting to import default from 'os' will result in mock value being incorrect
// ref: https://github.com/nodejs/node/blob/14f02fc2f7c1ea7989bdfeddfadc14921edd4e25/lib/os.js#L310
import { homedir } from 'os';
import path from 'path';

const homeDir = homedir();

const EDDY_DIR = path.join(homeDir, '.eddy.sh');
const EDDY_BIN_DIR = path.join(EDDY_DIR, 'bin');

export { EDDY_DIR, EDDY_BIN_DIR };
