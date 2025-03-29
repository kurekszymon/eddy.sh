#ifndef LANGUAGE_H
#define LANGUAGE_H

#include <functional>
#include <iostream>
#include <memory>
#include <string>
#include <unordered_map>
#include <vector>

#include "../ShellWrapper/ShellWrapper.hpp"

struct ToolInfo {
  // probably need reworking a bit
  bool loaded;
  std::string url;
  std::string version;
  std::function<void(const std::string &url)> install;
  const std::string placeholder = "{version}";

  ToolInfo(
      const std::string &url, const std::string &version, bool loaded,
      std::function<void(const std::string &url)> install =
          [](const std::string &url) { return 0; })

      : url(url), version(version), loaded(loaded), install(install) {}

  void set_version(const std::string &ver) {
    version = ver;
    format_url();
  }
  void set_loaded(bool load) { loaded = load; }

private:
  void format_url() {
    size_t pos = url.find(placeholder);
    while (pos != std::string::npos) {
      url.replace(pos, placeholder.length(), version);
      pos = url.find(placeholder, pos + version.length());
    }
  }
};

typedef std::unordered_map<std::string, ToolInfo> ToolMap;
typedef std::pair<std::string, ToolInfo> LoadedTool;

class Language {
public:
  virtual const ToolMap &get_tools() const = 0;
  virtual std::string get_name() const = 0;

  Language(std::shared_ptr<ShellWrapper> shell_) : shell(shell_) {}

  virtual void load_tool(const std::string &name, const std::string &version) {
    const ToolMap &tools = get_tools();

    for (const auto &tool : tools) {
      if (tool.first != name) {
        // add info about not loading tools
        continue;
      }

      ToolInfo tool_info = tool.second;
      tool_info.set_loaded(true);
      tool_info.set_version(version);

      this->loaded_tools.push_back({name, tool_info});
    }
  }

  virtual std::vector<LoadedTool> get_loaded_tools() { return loaded_tools; }

  virtual ~Language() = default;

protected:
  std::vector<LoadedTool> loaded_tools;
  std::shared_ptr<ShellWrapper> shell;
};

#endif // LANGUAGE_H
