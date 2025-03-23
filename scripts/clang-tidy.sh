#!/bin/zsh

clang-tidy -p conan $EDDY_PATH/src/**/*.{cpp,hpp}
clang-tidy -p conan $EDDY_PATH/src/main.cpp