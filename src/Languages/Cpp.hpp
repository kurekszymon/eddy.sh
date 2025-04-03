#ifndef CPP_H
#define CPP_H

#include "../Globals/Globals.hpp"
#include "Language.hpp"

namespace {

std::string get_ninja_platform() {
  // Base URL for Ninja releases
  std::string platform;
#ifdef _WIN32
  platform = "win";
#elif defined(__APPLE__)
  platform = "mac";
#else
  platform = "linux";
#endif

  return platform;
}

ArchiveType get_cmake_archive_enum() {
  ArchiveType extension_enum;

#ifdef _WIN32
  extension_enum = ArchiveType::ZIP;
#else
  extension_enum = ArchiveType::TAR_GZ;
#endif

  return extension_enum;
}
std::string get_cmake_archive_type() {
  std::string extension;

#ifdef _WIN32
  extension = ".zip";
#else
  extension = ".tar.gz";
#endif

  return extension;
}
} // namespace
class Cpp : public Language {
public:
  Cpp(std::shared_ptr<ShellWrapper> shell) : Language(shell) {}

  const ToolMap &get_tools() const override {
    static const ToolMap tools = {
        {"emscripten",
         ToolInfo("https://github.com/emscripten-core/emsdk.git", "latest",
                  false,
                  std::bind(&Cpp::install_emsdk, this, std::placeholders::_1,
                            std::placeholders::_2))},
        {"cmake",
         ToolInfo("https://github.com/Kitware/CMake/releases/download/"
                  "v{version}/cmake-{version}" +
                      get_cmake_archive_type(),
                  "latest", false,
                  std::bind(&Cpp::install_cmake, this, std::placeholders::_1,
                            std::placeholders::_2))},
        {"ninja",
         ToolInfo("https://github.com/ninja-build/ninja/releases/download/"
                  "v{version}/ninja-" +
                      get_ninja_platform() + ".zip",
                  "1.12.1", false,
                  std::bind(&Cpp::install_ninja, this, std::placeholders::_1,
                            std::placeholders::_2))}};
    return tools;
  }

  void install_cmake(const std::string &url, const std::string &version) const {
    shell->echo("$: CMake " + version + " installation starts");
    auto [dir, filename] = shell->curl(url, "cmake");
    auto extracted_filename =
        shell->extract(dir, filename, get_cmake_archive_enum());

    // create a function to concat dirs and filenames
    shell->run_script_file(dir + "/" + extracted_filename, "bootstrap");
    shell->run_make(dir + "/" + extracted_filename);
    // add to .bashrc PATH/cmake/bin
  }

  void install_emsdk(const std::string &url, const std::string &version) const {
    shell->echo("$: EMSDK " + version + " installation starts");
    auto dir = parse_home_dir(EDDY_PATH) + "/emsdk";

    shell->git_clone(url, dir);
    shell->git_pull(dir);

    shell->echo("# ./emsdk install " + version);
    shell->run_script_file(dir, "emsdk", "install", version);

    shell->echo("# ./emsdk activate latest");
    shell->run_script_file(dir, "emsdk", "activate", version);

    // blocking when script file run
    // echo 'source "~/.eddy.sh/emsdk/emsdk_env.sh"' >> $HOME/.zprofile
  }

  void install_ninja(const std::string &url, const std::string &version) const {
    shell->echo("$: Ninja " + version + " installation starts");

    auto [dir, filename] = shell->curl(url, "ninja");
    shell->extract(dir, filename, ArchiveType::ZIP);
    shell->make_executable(dir, "ninja");

    // move to ~/.eddy.sh/bin
  }
  std::string get_name() const override { return "Cpp"; }
};
#endif // CPP_H