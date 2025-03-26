
#ifndef SHELL_H
#define SHELL_H

#include <boost/asio.hpp>
#include <filesystem>
#include <iostream>
#include <memory>
#include <vector>

#include "boost/filesystem.hpp"
#include "boost/process.hpp"

#include "../Globals/Globals.hpp"
#include "../Log/Logger.hpp"
#include "../Yaml/CustomScript.hpp"

namespace bp = boost::process;
namespace bfs = boost::filesystem;
namespace fs = std::filesystem;

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

std::string check_eddy_path() {
  std::string eddy_path = parse_home_dir(EDDY_PATH);

  if (!std::filesystem::exists(eddy_path)) {
    std::filesystem::create_directory(eddy_path);
  }

  return parse_home_dir(eddy_path);
}
} // namespace

class Shell {
public:
  explicit Shell(std::shared_ptr<Logger> logger) : logger(logger) {}

  std::string pipe_output(bp::ipstream &pipe_stream, std::string src) {
    std::string output;

    std::string line;
    while (std::getline(pipe_stream, line)) {
      logger->update(line, src + ": ");
    }

    pipe_stream.close();
    return output;
  };

  void custom_command(CustomScript command) {
    bp::ipstream std_out;
    bp::ipstream std_err;

    std::vector<std::string> args = {"-c", command.cmd};

    bp::child c(bp::search_path("sh"), args, bp::std_out > std_out,
                bp::std_err > std_err);

    pipe_output(std_out, command.name);
    pipe_output(std_err, command.name);
    c.wait();

    if (c.exit_code() > 0) {
      std::string code = std::to_string(c.exit_code());
      logger->update(command.name + " exited with code: " + code);
    }
  }

  void run_script_file(const std::string &path) {
    bp::ipstream std_out;

    std::vector<std::string> args = {"-c", parse_home_dir(path)};

    bp::child c(bp::search_path("sh"), args, bp::std_out > std_out);

    int bash_exit = c.exit_code();

    pipe_output(std_out, path);

    c.wait();

    if (c.exit_code() > 0) {
      std::string code = std::to_string(c.exit_code());
      logger->update(path + " exited with code: " + code);
    }
  }

  void make_executable(const std::string &filepath) {
    auto _filepath = parse_home_dir(filepath);
    logger->update("chmod +x " + _filepath);

    auto perms =
        fs::perms::owner_exec | fs::perms::group_exec | fs::perms::others_exec;

    fs::permissions(_filepath, perms, fs::perm_options::add);

    logger->update("File made executable: " + _filepath);
  };
  void git_clone(const std::string &repo_url, const std::string &clone_dir) {
    // TODO: handle git clone async command in logger
    bp::ipstream std_err;
    std::vector<std::string> args = {"clone", repo_url,
                                     parse_home_dir(clone_dir)};

    logger->update("$: Starting cloning " + repo_url + " to " + clone_dir);
    bp::child c(bp::search_path("git"), args,
                bp::std_err > std_err); // git clone writes to std_err.

    c.wait();
    pipe_output(std_err, "git clone");

    if (c.exit_code() > 0) {
      std::string code = std::to_string(c.exit_code());
      logger->update("$: Cloning failed with code " + code);
    } else {
      logger->update("$: Successfuly cloned " + repo_url + ".");
    }
  };

  enum curl_type { download };

  std::string curl(const std::string url, const std::string name,
                   curl_type type = curl_type::download) {
    const std::string eddy_path = check_eddy_path();
    const std::string output_dir = eddy_path + "/" + name + ".sh";
    const std::vector<std::string> args = {
        url,
        "-o",
        output_dir,
    };
    bp::ipstream std_out;
    bp::ipstream std_err;
    bp::child c(bp::search_path("curl"), args, bp::std_out > std_out,
                bp::std_err > std_err);

    c.wait();
    pipe_output(std_out, "curl");
    pipe_output(std_err, "curl");

    if (c.exit_code() > 0) {
      logger->update("curl failed with code " + std::to_string(c.exit_code()));
    } else {
      logger->update("download complete.", "curl: ");
    }

    return output_dir;
  }

  void echo(const std::string msg) { logger->update(msg); }

private:
  std::shared_ptr<Logger> logger;
};
#endif // SHELL_H
