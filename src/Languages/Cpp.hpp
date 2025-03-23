#ifndef CPP_H
#define CPP_H

#include "Language.hpp"

class Cpp : public Language {
public:
  Cpp(std::shared_ptr<Shell> shell) : Language(shell) {}

  const ToolMap &get_tools() const override {
    static const ToolMap tools = {
        {"emscripten", ToolInfo("url", "latest", false)},
        {"cmake", ToolInfo("ur;", "latest", false)},
        {"ninja", ToolInfo("url;", "3", false)}};
    return tools;
  }

  std::string get_name() const override { return "Cpp"; }
};
#endif // CPP_H