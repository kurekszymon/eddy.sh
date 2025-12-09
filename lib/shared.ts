import fs from 'fs';
import path from 'path';

import { EDDY_DIR } from './consts';

export const createToolDir = (dirName: string) => {
    return fs.mkdirSync(path.join(EDDY_DIR, dirName), { mode: 0o755, recursive: true });
};
