#ifndef CPP_H
#define CPP_H

#include "../Globals/Globals.hpp"
#include "Language.hpp"

namespace {
ArchiveType get_compression_extension_enum() {
  ArchiveType extension_enum;

#ifdef _WIN32
  extension_enum = ArchiveType::ZIP;
#else
  extension_enum = ArchiveType::TAR_GZ;
#endif

  return extension_enum;
}
std::string get_compression_extension() {
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
        {"emscripten", ToolInfo("url", "latest", false)},
        {"cmake",
         ToolInfo("https://github.com/Kitware/CMake/releases/download/"
                  "v{version}/cmake-{version}" +
                      get_compression_extension(),
                  "latest", false,
                  std::bind(&Cpp::install_cmake, this, std::placeholders::_1))},
        {"ninja", ToolInfo("url;", "3", false)}};
    return tools;
  }

  void install_cmake(const std::string &url) const {
    shell->echo("installing cmake..");
    shell->echo(url);
    auto [dir, filename] = shell->curl(url, "cmake");
    auto extracted_filename =
        shell->tar(dir, filename, get_compression_extension_enum());

    shell->bootstrap_cmake(extracted_filename);
  }

  std::string get_name() const override { return "Cpp"; }
};
#endif // CPP_H