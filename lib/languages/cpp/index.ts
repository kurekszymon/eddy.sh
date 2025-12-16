import { bazel } from './bazel';
import { cmake } from './cmake';
import { conan } from './conan';
import { ninja } from './ninja';

export const cpp = {
    emscripten: () => {
        return 'hello from emscripten';
    },
    bazel,
    cmake,
    conan,
    ninja,
} as const;