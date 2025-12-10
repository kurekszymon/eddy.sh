import { ninja } from "./ninja";

export const cpp = {
    cmake: () => {
        return 'hello from cmake';
    },
    ninja,
    emscripten: () => {
        return 'hello from emscripten';
    },
} as const;