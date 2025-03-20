#!/bin/zsh

clang-tidy -p conan $EDDY_PATH/src/*/**.hpp
clang-tidy -p conan $EDDY_PATH/src/*.{cpp,hpp}