#ifndef JAVASCRIPT_H
#define JAVASCRIPT_H

#include "Language.hpp"

class Javascript : public Language {
public:
  Javascript(std::shared_ptr<ShellWrapper> shell) : Language(shell) {}
  const ToolMap &get_tools() const override {
    static const ToolMap tools = {
        {"nvm", ToolInfo("https://raw.githubusercontent.com/nvm-sh/nvm/"
                         "v{version}/install.sh",
                         "latest", false,
                         std::bind(&Javascript::install_nvm, this,
                                   std::placeholders::_1))}};
    return tools;
  }

  void install_nvm(const std::string &url) const {
    shell->echo("Downloading nvm...");

    auto [output_dir, filename] = shell->curl(url, "nvm");
    shell->make_executable(output_dir, filename);
    shell->run_script_file(output_dir, filename);
  }
  std::string get_name() const override { return "Javascript"; }

private:
  const ToolMap tools;
};
#endif // JAVASCRIPT_H