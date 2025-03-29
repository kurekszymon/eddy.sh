#ifndef PYTHON_H
#define PYTHON_H

#include "Language.hpp"

class Python : public Language {
public:
  Python(std::shared_ptr<ShellWrapper> shell) : Language(shell) {}

  const ToolMap &get_tools() const override {
    static const ToolMap tools = {{"pyenv", ToolInfo("ur;", "latest", false)},
                                  {"pip", ToolInfo("ur;'", "3", false)}};
    return tools;
  }

  std::string get_name() const override { return "Python"; }
};
#endif // PYTHON_H