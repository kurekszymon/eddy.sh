
#include <boost/asio.hpp>
#include <iostream>
#include <memory>
#include <vector>

#include "boost/filesystem.hpp"
#include "boost/process.hpp"

#include "../Config.hpp" // only for CustomScript type import, move it to seperate file
#include "../Log/Logger.hpp"
namespace bp = boost::process;
namespace bfs = boost::filesystem;

namespace {
std::string parse_home_dir(std::string input_path) {
  if (input_path.empty() || input_path[0] != '~') {
    return input_path;
  }

  const char *env_home = std::getenv("HOME"); // Unix-like systems (Linux/macOS)
  const char *env_userprofile = std::getenv("USERPROFILE"); // Windows

  bfs::path home_dir = env_home != nullptr ? env_home : env_userprofile;

  // handle home not found - decide for try catch / throwing

  bfs::path result = home_dir / input_path.substr(1);
  return result.string();
};
} // namespace

class Shell {
public:
  explicit Shell(std::shared_ptr<Logger> logger) : logger_(logger) {}

  std::string pipe_output(bp::ipstream &pipe_stream, std::string src) {
    std::string output;

    std::string line;
    while (std::getline(pipe_stream, line)) {
      logger_->update(src + ": " + line);
    }

    pipe_stream.close();
    return output;
  };

  std::string execute_custom_command(CustomScript command) {
    bp::ipstream std_out;
    bp::ipstream std_err;

    std::vector<std::string> args = {"-c", command.cmd};

    bp::child c(bp::search_path("sh"), args, bp::std_out > std_out,
                bp::std_err > std_err);

    std::string output = pipe_output(std_out, command.name);
    std::string err_output = pipe_output(std_err, command.name);

    c.wait();

    if (c.exit_code() > 0) {
      std::string code = std::to_string(c.exit_code());
      if (!err_output.empty()) {
        return command.name + " exited with code: " + code + " " + err_output;
      }

      return command.name + " exited with code: " + code;
    }

    return output;
  }

  std::string execute_git_clone(const std::string &repo_url,
                                const std::string &clone_dir) {
    // TODO: handle git clone async command in logger
    bp::ipstream std_err;
    std::vector<std::string> args = {"clone", repo_url,
                                     parse_home_dir(clone_dir)};

    bp::child c(bp::search_path("git"), args,
                bp::std_err > std_err); // git clone writes to std_err.

    std::string output = pipe_output(std_err, "git clone");

    c.wait();

    return output;
  };

private:
  std::shared_ptr<Logger> logger_;
};
