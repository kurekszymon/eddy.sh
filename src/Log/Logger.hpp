#ifndef LOGGER_H
#define LOGGER_H

#include <ftxui/component/component.hpp>
#include <ftxui/dom/elements.hpp>
#include <string>
#include <vector>

using namespace ftxui;

namespace {
Elements AddFocusBottom(Elements list) {
  if (!list.empty()) {
    // yframe needs to have focus to know what element to display
    // this ensures focus is at the end of the list.
    list.back() = focus(std::move(list.back()));
  }
  return std::move(list);
}
} // namespace

class Logger {
public:
  void update(const std::string &new_line, std::string prefix = "$: ") {
    lines.push_back(new_line);
  }

  Component renderer = Renderer([this] {
    return vbox({this->render_console_output()}) | yframe | yflex | border;
  });

private:
  std::vector<std::string> lines;
  std::function<Element()> render_console_output = [this]() {
    std::vector<Element> elements;

    for (const auto &line : this->lines) {
      elements.push_back(paragraph(line));
    }

    if (lines.empty()) {
      return vbox(text("Execute some command to see the output") | center);
    }

    return vbox(AddFocusBottom(elements));
  };
};
#endif // LOGGER_H
