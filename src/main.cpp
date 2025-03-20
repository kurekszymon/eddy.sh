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

#include "Config.hpp"
#include "Log/Logger.hpp"
#include "Shell/Shell.hpp"

using namespace ftxui;

int main() {
  auto logger = std::make_shared<Logger>();
  auto shell = std::make_unique<Shell>(logger);
  auto config = std::make_unique<Config>("config.yaml");
  auto screen = ScreenInteractive::Fullscreen();

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

    shell->execute_custom_command(selected_script);
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

    shell->execute_git_clone(selected_repo.url, clone_to);
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

    std::string console_output;
    for (const auto &tool : tools) {
      console_output.append(tool.first + " ");
    }

    auto render_console_output = [console_output]() {
      return hbox({text("Lang console output: "), text(console_output)}) |
             border | flex;
    };
    auto lang_renderer = Renderer([render_console_output] {
      return vbox({render_console_output()}) | flex;
    });

    auto lang_container = Container::Vertical({lang_renderer});

    if (!tools.empty()) {
      tab_entries.push_back(language->get_name());
      tab_content->Add(lang_container);
    }
  }

  // ---------------------------------------------------------------------------
  // Main render
  // ---------------------------------------------------------------------------

  auto tab_selection =
      Menu(&tab_entries, &tab_index, MenuOption::HorizontalAnimated());

  auto exit_button =
      Button("Exit", [&] { screen.Exit(); }, ButtonOption::Animated());

  auto logger_renderer = logger->renderer;

  auto tab_selection_container = Container::Horizontal({
      tab_selection,
      exit_button,
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
  auto group_renderer = ResizableSplitLeft(main_renderer, logger_renderer,
                                           &paragraph_renderer_split_position);

  std::atomic<bool> refresh_ui_continue = true;
  std::thread refresh_ui([&] {
    while (refresh_ui_continue) {
      using namespace std::chrono_literals;
      std::this_thread::sleep_for(0.05s);
      // The |shift| variable belong to the main thread. `screen.Post(task)`
      // will execute the update on the thread where |screen| lives (e.g. the
      // main thread). Using `screen.Post(task)` is threadsafe.
      screen.Post([&] { shift++; });
      // After updating the state, request a new frame to be drawn. This is done
      // by simulating a new "custom" event to be handled.
      screen.Post(Event::Custom);
    }
  });

  screen.Loop(group_renderer);
  refresh_ui_continue = false;
  refresh_ui.join();

  return 0;
}