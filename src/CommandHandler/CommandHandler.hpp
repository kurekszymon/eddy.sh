#include <atomic>
#include <boost/filesystem.hpp>
#include <boost/process.hpp>
#include <ftxui/component/component.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <ftxui/dom/elements.hpp>
#include <mutex>
#include <string>
#include <thread>
#include <vector>

namespace bp = boost::process;
namespace bfs = boost::filesystem;

class CommandHandler {
private:
  std::vector<std::string> lines;
  std::mutex lines_mutex;
  std::mutex command_mutex;
  std::atomic<bool> command_running{false};
  ftxui::ScreenInteractive &screen;

  void add_line(const std::string &line) {
    std::lock_guard<std::mutex> lock(lines_mutex);
    lines.push_back(line);
  }

public:
  CommandHandler(ftxui::ScreenInteractive &scr) : screen(scr) {}

  ftxui::Element render_console_output() {
    std::lock_guard<std::mutex> lock(lines_mutex);
    std::vector<ftxui::Element> elements;

    for (const auto &line : lines) {
      elements.push_back(ftxui::paragraph(line) | ftxui::dim);
    }

    if (command_running) {
      elements.push_back(ftxui::text(L"Wait for command to finish..") |
                         ftxui::focus);
    } else {
      elements.push_back(ftxui::text(L"Command finished!") | ftxui::focus);
    }

    if (lines.empty()) {
      return ftxui::vbox(ftxui::text("Execute a command to see the output") |
                         ftxui::center);
    }

    return ftxui::vbox(elements) | ftxui::yframe;
  }

  void
  run_command(const std::string &command,
              const std::string &working_dir = bfs::current_path().string()) {

    while (command_running) {
      std::this_thread::sleep_for(std::chrono::milliseconds(50));
    }
    command_running = true;
    std::thread command_thread([this, command, working_dir]() {
      bp::ipstream std_out;
      bp::ipstream std_err;
      bp::child c(command, bp::start_dir(working_dir), bp::std_out > std_out,
                  bp::std_err > std_err);

      std::string line;
      while (std::getline(std_out, line)) {
        // Add line and request screen update
        screen.Post([this, line]() {
          add_line(line);
          screen.Post(ftxui::Event::Custom);
        });
      }

      while (std::getline(std_err, line)) {
        // Add line and request screen update
        screen.Post([this, line]() {
          add_line(line);
          screen.Post(ftxui::Event::Custom);
        });
      }

      // Wait for process to complete
      c.wait();
      screen.Post(ftxui::Event::Custom); // post custom event to refresh ui
      // after command finishes
      command_running = false;
    });

    // Detach the thread so it can run independently -> dont block the ui
    command_thread.detach();
  }

  bool is_command_running() const { return command_running; }
};
