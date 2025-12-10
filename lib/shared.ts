import fs from 'fs';
import path from 'path';
import https from 'https';

import { EDDY_DIR } from './consts';
import { logger } from './logger';

export const createToolDir = (dirName: string) => {
    const dir = path.join(EDDY_DIR, dirName);

    try {
        fs.mkdirSync(dir, { mode: 0o755, recursive: true });
    } catch (e) {
        // TODO: handle error
        logger.error(e);
    }

    return dir;
};


export const downloadFile = (filePath: string, url: string) => {
    const file = fs.createWriteStream(filePath);

    return new Promise<void>((res, rej) => {
        https.get(url, (response) => {
            const total = parseInt(response.headers['content-length'] || '0', 10);
            let downloaded = 0;

            response.on('data', (chunk) => {
                downloaded += chunk.length;
                if (total) {
                    const percent = ((downloaded / total) * 100).toFixed(2);
                    logger.info(`Downloading progress: ${percent}%`);
                } else {
                    logger.info(`Downloading progress: ${downloaded} bytes`);
                }
            });

            response.pipe(file);

            // after download completed close filestream
            file.on("finish", () => {
                file.close();
                logger.info(`Download complete!`);
                res();
            });
        }).on('error', (err) => {
            // TODO: handle error
            logger.error('http request error');
            rej(err);
        });

        file.on('error', (err) => {
            // TODO: handle error
            logger.error('file write error');
            rej(err);
        });
    });
};