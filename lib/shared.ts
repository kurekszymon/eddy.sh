import fs from 'fs';
import { EDDY_BIN_DIR } from './consts';

export const createEddyDirs = () => {
    return fs.mkdirSync(EDDY_BIN_DIR, { mode: 0o755, recursive: true });
};
