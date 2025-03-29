#include <array>  // for array
#include <atomic> // for atomic
#include <chrono> // for operator""s, chrono_literals
#include <iostream>
#include <memory> // for make_unique
#include <string> // for string, basic_string, char_traits, operator+, to_string
#include <thread> // for sleep_for, thread
#include <vector> // for vector

#include "ftxui/component/component.hpp" // for Checkbox, Renderer, Horizontal, Vertical, Input, Menu, Radiobox, ResizableSplitLeft, Tab
#include "ftxui/component/component_base.hpp"    // for ComponentBase, Component
#include "ftxui/component/component_options.hpp" // for MenuOption, InputOption
#include "ftxui/component/event.hpp"             // for Event, Event::Custom
#include "ftxui/component/screen_interactive.hpp" // for Component, ScreenInteractive
#include "ftxui/dom/elements.hpp" // for text, color, operator|, bgcolor, filler, Element, vbox, size, hbox, separator, flex, window, graph, EQUAL, paragraph, WIDTH, hcenter, Elements, bold, vscroll_indicator, HEIGHT, flexbox, hflow, border, frame, flex_grow, gauge, paragraphAlignCenter, paragraphAlignJustify, paragraphAlignLeft, paragraphAlignRight, dim, spinner, LESS_THAN, center, yframe, GREATER_THAN

#include "Config/Config.hpp"
#include "ShellWrapper/ShellWrapper.hpp"

using namespace ftxui;

int main() {
  auto screen = ScreenInteractive::Fullscreen();
  auto command_handler = std::make_shared<CommandHandler>(screen);
  auto shell = std::make_shared<ShellWrapper>(command_handler);
  auto config = std::make_unique<Config>(shell);

  int shift = 0;
  int tab_index = 0;
  std::vector<std::string> tab_entries = {};

  auto tab_content = Container::Tab({}, &tab_index);

  // TODO move sections to seperate functions

  // ---------------------------------------------------------------------------
  // Custom Scripts
  // ---------------------------------------------------------------------------

  auto custom_scripts_config = config->custom_scripts;
  int cs_selected_index = 0;

  std::vector<std::string> cs_entries;
  for (const auto &script : custom_scripts_config) {
    cs_entries.push_back(script.name);
  }

  auto cs_option = MenuOption::VerticalAnimated();

  cs_option.on_enter = [&]() {
    const CustomScript &selected_script =
        custom_scripts_config.at(cs_selected_index);

    shell->run_custom_command(selected_script);
  };

  auto cs_menu = Menu(&cs_entries, &cs_selected_index, cs_option);

  auto cs_container = Container::Horizontal({cs_menu});

  if (!config->custom_scripts.empty()) {
    tab_entries.push_back("custom scripts");
    tab_content->Add(cs_container);
  }
  // ---------------------------------------------------------------------------
  // Repositories
  // ---------------------------------------------------------------------------

  auto repo_config = config->repositories;
  int repositories_selected_index = 0;

  std::vector<std::string> repositories_entries;
  for (const auto &repo : repo_config.vector) {
    repositories_entries.push_back(repo.name);
  }

  auto repositories_option = MenuOption::VerticalAnimated();
  repositories_option.on_enter = [&]() {
    const Repository &selected_repo =
        repo_config.vector.at(repositories_selected_index);
    const std::string clone_to =
        repo_config.clone_path + '/' + selected_repo.name;

    shell->git_clone(selected_repo.url, clone_to);
  };

  auto repos_menu = Menu(&repositories_entries, &repositories_selected_index,
                         repositories_option);

  auto repositories_container = Container::Horizontal({repos_menu});

  if (!config->repositories.vector.empty()) {
    tab_entries.push_back("repositories");
    tab_content->Add(repositories_container);
  }

  // ---------------------------------------------------------------------------
  // Languages
  // ---------------------------------------------------------------------------

  auto languages_config = config->languages;
  for (const auto &language : languages_config) {
    auto tools = language->get_loaded_tools();
    if (tools.empty()) {
      continue;
    }
    struct SharedState {
      // needed to properly handle lifetimes
      std::vector<std::string> entries;
      std::vector<LoadedTool> tool_data;
      int selected = 0;
    };

    auto state = std::make_shared<SharedState>();

    for (const auto &tool : tools) {
      state->entries.push_back(tool.first);
      state->tool_data.push_back(tool);
    }

    auto menu = Menu(&state->entries, &state->selected);

    auto menu_with_event = CatchEvent(menu, [state](const Event &event) {
      if (event == Event::Return) {
        if (state->selected >= 0 &&
            state->selected < static_cast<int>(state->entries.size())) {

          const auto &selected_tool = state->tool_data[state->selected];

          const auto &tool_info = selected_tool.second;

          tool_info.install(tool_info.url);
        }

        return true;
      }
      return false;
    });

    Component renderer = Renderer(menu_with_event, [menu_with_event, state] {
      return menu_with_event->Render() | yframe | yflex | border;
    });

    tab_entries.push_back(language->get_name());
    tab_content->Add(renderer);
  }

  // ---------------------------------------------------------------------------
  // Main render
  // ---------------------------------------------------------------------------

  auto tab_selection =
      Menu(&tab_entries, &tab_index, MenuOption::HorizontalAnimated());

  auto exit_button =
      Button("Exit", [&] { screen.Exit(); }, ButtonOption::Animated());

  auto tab_selection_container = Container::Horizontal({
      tab_selection,
      exit_button,
  });

  auto console_renderer = ftxui::Renderer([&command_handler]() {
    return ftxui::vbox({command_handler->render_console_output()}) |
           ftxui::yframe | ftxui::yflex | ftxui::border;
  });

  auto main_container =
      Container::Vertical({tab_selection_container, tab_content});

  auto main_renderer = Renderer(main_container, [&] {
    return vbox({
        text("eddy.sh") | bold | hcenter,
        hbox({
            tab_selection->Render() | flex,
            exit_button->Render(),
        }),
        hbox({
            tab_content->Render(),
        }),
    });
  });

  int paragraph_renderer_split_position = Terminal::Size().dimx / 1.5;
  auto group_renderer = ResizableSplitLeft(main_renderer, console_renderer,
                                           &paragraph_renderer_split_position);

  screen.Loop(group_renderer);

  return 0;
}