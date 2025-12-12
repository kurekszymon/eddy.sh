import { cmake } from './cmake';
import { ninja } from './ninja';

export const cpp = {
    cmake,
    ninja,
    emscripten: () => {
        return 'hello from emscripten';
    },
} as const;