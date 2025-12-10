export const logger = {
    debug: (...msg: any) => console.log('[DEBUG]: ', ...msg),
    info: (...msg: any) => console.log('[INFO]: ', ...msg),
    warn: (...msg: any) => console.log('[WARN]: ', ...msg),
    error: (...msg: any) => console.log('[ERROR]: ', ...msg)
};