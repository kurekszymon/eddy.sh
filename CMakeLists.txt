cmake_minimum_required(VERSION 4.0.0)
project(eddy.sh VERSION 0.1)
include(FetchContent)

find_package(yaml-cpp REQUIRED)
find_package(Boost REQUIRED)

# use FetchContent to download ftxui - conan has outdated version (needed for Ctrl+C support)
FetchContent_Declare(ftxui
  GIT_REPOSITORY https://github.com/ArthurSonzogni/ftxui
  GIT_TAG v6.0.2
)

FetchContent_GetProperties(ftxui)
if(NOT ftxui_POPULATED)
  FetchContent_MakeAvailable(ftxui)
endif()

find_program(CLANG_TIDY_EXECUTABLE NAMES clang-tidy)

if(CLANG_TIDY_EXECUTABLE)
  set(CMAKE_CXX_CLANG_TIDY ${CLANG_TIDY_EXECUTABLE})
else()
  message(
    STATUS "clang-tidy not found, skipping CMake CXX_CLANG_TIDY configuration.")
endif()

add_subdirectory(src)

include_directories(${YAMLCPP_INCLUDE_DIRS})
include_directories(${CMAKE_SOURCE_DIR}/include)
include_directories(${Boost_INCLUDE_DIR})

target_link_libraries(
  eddy.sh

  ftxui::screen
  ftxui::dom
  ftxui::component

  boost::boost
  yaml-cpp::yaml-cpp
)


