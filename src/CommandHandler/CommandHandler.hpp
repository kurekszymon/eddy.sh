#include <atomic>
#include <mutex>
#include <string>
#include <thread>
#include <vector>

#include <boost/filesystem.hpp>
#include <boost/process.hpp>
#include <ftxui/component/component.hpp>
#include <ftxui/component/screen_interactive.hpp>
#include <ftxui/dom/elements.hpp>

namespace bp = boost::process;
namespace bfs = boost::filesystem;

class CommandHandler {
private:
  std::vector<std::string> lines;
  std::mutex lines_mutex;
  std::atomic<bool> command_running{false};
  ftxui::ScreenInteractive &screen;
  std::queue<std::pair<std::string, std::string>> command_queue;
  std::mutex queue_mutex;
  std::condition_variable cv;
  std::thread worker_thread;
  std::atomic<bool> shutdown{false};

  void add_line(const std::string &line) {
    std::lock_guard<std::mutex> lock(lines_mutex);
    lines.push_back(line);
  }

public:
  CommandHandler(ftxui::ScreenInteractive &scr) : screen(scr) {
    worker_thread = std::thread(&CommandHandler::process_command_queue, this);
  }

  ~CommandHandler() {
    shutdown = true;
    cv.notify_all();
    if (worker_thread.joinable()) {
      worker_thread.join();
    }
  }

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
    std::lock_guard<std::mutex> lock(queue_mutex);
    command_queue.push({command, working_dir});
    cv.notify_one();
  }

  bool is_command_running() const { return command_running; }

private:
  // Worker thread function that processes the command queue sequentially
  void process_command_queue() {
    while (!shutdown) {
      std::unique_lock<std::mutex> lock(queue_mutex);
      cv.wait(lock, [this]() { return !command_queue.empty() || shutdown; });

      if (shutdown) {
        break;
      }

      if (!command_queue.empty()) {
        auto cmd_pair = command_queue.front();
        command_queue.pop();
        lock.unlock();

        execute_command(cmd_pair.first, cmd_pair.second);
      }
    }
  }

  void execute_command(const std::string &command,
                       const std::string &working_dir) {
    command_running = true;
    screen.Post(ftxui::Event::Custom);

    bp::ipstream std_out;
    bp::ipstream std_err;
    bp::child c(command, bp::start_dir(working_dir), bp::std_out > std_out,
                bp::std_err > std_err);

    std::string line;
    while (std::getline(std_out, line)) {
      screen.Post([this, line]() {
        add_line(line);
        screen.Post(ftxui::Event::Custom);
      });
    }

    while (std::getline(std_err, line)) {
      screen.Post([this, line]() {
        add_line(line);
        screen.Post(ftxui::Event::Custom);
      });
    }

    c.wait();

    command_running = false;
    screen.Post(ftxui::Event::Custom);
  }
};