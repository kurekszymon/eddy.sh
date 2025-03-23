#ifndef JAVASCRIPT_H
#define JAVASCRIPT_H

#include "Language.hpp"

class Javascript : public Language {
public:
  Javascript(std::shared_ptr<Shell> shell) : Language(shell) {}

  const ToolMap &get_tools() const override {
    static const ToolMap tools = {
        {"nvm",
         ToolInfo("https://raw.githubusercontent.com/nvm-sh/nvm/{version}/"
                  "install.sh",
                  "latest", false, install_nvm)},
    };
    return tools;
  }

  std::string get_name() const override { return "Javascript"; }
  // learn how to replace this lambda with something else.
  std::function<int()> install_nvm = [this]() {
    this->shell_->echo("hello from install nvm");
    return 1;
  };

private:
  const ToolMap tools;
};
#endif // JAVASCRIPT_H