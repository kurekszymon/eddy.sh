#include <array>
#include <chrono>
#include <iostream>

#include "ftxui/component/component.hpp"
#include "ftxui/component/component_base.hpp"
#include "ftxui/component/component_options.hpp"
#include "ftxui/component/event.hpp"
#include "ftxui/component/screen_interactive.hpp"
#include "ftxui/dom/elements.hpp"

#include "Config/Config.hpp"
#include "ShellWrapper/ShellWrapper.hpp"

using namespace ftxui;

#define LOGGER_RATIO 1.75

int main() {
  auto screen = ScreenInteractive::Fullscreen();
  // Catch Ctrl+C to add ~/.eddy.sh to shell, then quit.
  screen.ForceHandleCtrlC(false);

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

    auto menu_with_event = CatchEvent(menu, [state, &command_handler](
                                                const Event &event) {
      if (event == Event::Return && !command_handler->is_command_running()) {
        if (state->selected >= 0 &&
            state->selected < static_cast<int>(state->entries.size())) {

          const auto &selected_tool = state->tool_data[state->selected];

          const auto &tool_info = selected_tool.second;

          tool_info.install(tool_info.url, tool_info.version);
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

  auto exit_button = Button(
      "Exit",
      [&] {
        auto used_shell = detect_used_shell();
        if (!used_shell) {
          std::cout << "Unsupported shell type. Please add ~/.eddy.sh/bin to "
                       "your .shellrc";
        }

        std::vector<std::string> command = {
            "echo", "'export PATH=\"~/.eddy.sh/bin:$PATH\"'", ">> ~/.zshrc"};
        command_handler->run_command(command);
        screen.Exit();
      },
      ButtonOption::Animated());

  auto tab_selection_container = Container::Horizontal({
      tab_selection,
      exit_button,
  });

  auto console_renderer = ftxui::Renderer([&command_handler]() {
    return ftxui::vbox({command_handler->render_console_output()}) |
           ftxui::yframe | ftxui::yflex | ftxui::border;
  });

  auto _main_container =
      Container::Vertical({tab_selection_container, tab_content});

  auto main_container = CatchEvent(
      _main_container, [&command_handler, &shell, &screen](const Event &event) {
        if (event == Event::CtrlC) {

          auto used_shell = detect_used_shell();
          if (!used_shell) {
            std::cout << "Unsupported shell type. Please add ~/.eddy.sh/bin to "
                         "your .shellrc";
          }

          std::vector<std::string> command = {
              "echo", "'export PATH=\"~/.eddy.sh/bin:$PATH\"'", ">> ~/.zshrc"};
          command_handler->run_command(command);
          screen.Exit();

          return true;
        }
        return false;
      });

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

  int paragraph_renderer_split_position = Terminal::Size().dimx / LOGGER_RATIO;
  auto group_renderer = ResizableSplitLeft(main_renderer, console_renderer,
                                           &paragraph_renderer_split_position);

  screen.Loop(group_renderer);

  return 0;
}