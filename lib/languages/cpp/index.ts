export const cpp = {
    cmake: () => {
        return 'hello from cmake';
    },
    ninja: () => {
        return 'hello from ninja';
    },
    emscripten: () => {
        return 'hello from emscripten';
    },
} as const;