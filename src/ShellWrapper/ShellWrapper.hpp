
#ifndef SHELL_WRAPPER_H
#define SHELL_WRAPPER_H

#include <filesystem>
#include <memory>
#include <optional>
#include <vector>

#include <boost/algorithm/string/join.hpp>
#include <boost/asio.hpp>
#include <boost/filesystem.hpp>
#include <boost/process.hpp>

#include "../CommandHandler/CommandHandler.hpp"
#include "../Globals/Globals.hpp"
#include "../Yaml/CustomScript.hpp"

namespace bp = boost::process;
namespace bfs = boost::filesystem;
namespace fs = std::filesystem;

// auto *d = getenv("SHELL"); // add parsing
// refactor the returns to support exit codes
// std::pair<int, std::map<T>> probably?

namespace {
// since moving to boost::shell, probably not needed, get rid of
std::string parse_home_dir(const std::string input_path) {
  if (input_path.empty() || input_path[0] != '~') {
    return input_path;
  }

  const char *env_home =
      std::getenv("HOME"); // Unix-like systems (Linux / macOS)
  const char *env_userprofile = std::getenv("USERPROFILE"); // Windows

  bfs::path home_dir = env_home != nullptr ? env_home : env_userprofile;

  // handle home not found - decide for try catch / throwing

  bfs::path result = home_dir / input_path.substr(1);
  return result.string();
};

struct ShellConfig {
  std::string name;
  std::string rc;
  std::string profile;
};

std::optional<ShellConfig> detect_used_shell() {
  const char *shell = std::getenv("SHELL");
  std::string shellPath = shell ? std::string(shell) : "";

  if (shellPath.find("bash") != std::string::npos) {
    return ShellConfig{"bash", "~/.bashrc", "~/.bash_profile"};
  }

  if (shellPath.find("zsh") != std::string::npos) {
    return ShellConfig{"zsh", "~/.zshrc", "~/.zprofile"};
  }

  return std::nullopt;
};

std::string get_curl_filename(const std::string &url) {
  size_t lastSlashPos = url.find_last_of('/');
  std::string filename = url.substr(lastSlashPos + 1);

  return filename;
}

std::string get_extracted_filename(const std::string &filename) {
  std::vector<std::string> extensions = {".tar.gz", ".tar", ".zip"};

  for (const auto &ext : extensions) {
    if (filename.length() > ext.length() &&
        filename.substr(filename.length() - ext.length()) == ext) {
      return filename.substr(0, filename.length() - ext.length());
    }
  }

  return filename;
}

} // namespace

class ShellWrapper {
private:
  std::shared_ptr<CommandHandler> command_handler;

public:
  enum curl_type { download };
  std::string eddy_path = parse_home_dir(EDDY_PATH);
  std::string eddy_path_bin = eddy_path + "/bin";

  explicit ShellWrapper(std::shared_ptr<CommandHandler> handler)
      : command_handler(handler) {}

  std::string check_eddy_path() {
    std::string eddy_path = parse_home_dir(EDDY_PATH);

    if (!fs::exists(eddy_path)) {
      this->echo("$: ~/.eddy.sh doesn't exist, creating..");
      this->mkdir(eddy_path);
    }

    return eddy_path;
  }

  std::string check_eddy_bin_path() {
    std::string eddy_bin_path = parse_home_dir(eddy_path_bin);

    if (!fs::exists(eddy_bin_path)) {
      this->echo("$: ~/.eddy.sh/bin doesn't exist, creating..");
      this->mkdir(eddy_bin_path);
    }

    return eddy_bin_path;
  }

  void echo(const std::string &message) {
    std::vector<std::string> construct = {"echo", message};

    command_handler->run_command(construct);
  }

  void mkdir(const std::string &path) {
    std::vector<std::string> construct = {"mkdir", "-p", path};

    command_handler->run_command(construct);
  }

  void run_custom_command(CustomScript command_) {
    this->echo("$ " + command_.name + " started..");
    std::vector<std::string> construct = {command_.cmd};

    command_handler->run_command(construct);
  }

  std::pair<std::string, std::string>
  curl(const std::string url, const std::string name,
       curl_type type = curl_type::download) {

    const std::string eddy_path = check_eddy_path();
    const std::string curled_filename = get_curl_filename(url);
    const std::string output_path = eddy_path + "/" + curled_filename;

    const std::vector<std::string> construct = {
        "curl", "-L", "--output-dir", eddy_path, "-O", url,
    };

    command_handler->run_command(construct);

    return std::make_pair(eddy_path, curled_filename);
  }

  void make_executable(const std::string &file_dir,
                       const std::string &filename) {
    const std::string filepath = file_dir + "/" + filename;
    this->echo("$ chmod +x " + filepath);

    // causing issues because not scheduled
    std::vector<std::string> construct = {"chmod", "+x", filename};

    command_handler->run_command(construct, file_dir);
  };

  template <typename... Arguments>
  void run_script_file(const std::string &file_dir, const std::string &filename,
                       Arguments... arguments) {
    std::vector<std::string> construct = {"./" + filename};

    (construct.push_back(arguments), ...);

    command_handler->run_command(construct, file_dir);
  }

  void run_make(const std::string &dir) {
    std::vector<std::string> construct = {"make"};

    command_handler->run_command(construct, dir);
  }

  std::string extract(const std::string &path, const std::string name,
                      ArchiveType type, std::string extract_path = "") {
    std::string archive_path = path + "/" + name;

    if (extract_path.empty()) {
      extract_path = check_eddy_path();
    }

    if (type == ArchiveType::TAR_GZ) {
      std::vector<std::string> construct = {"tar", "-xzf", archive_path, "-C",
                                            extract_path};
      command_handler->run_command(construct);
    } else if (type == ArchiveType::TAR) {
      std::vector<std::string> construct = {"tar", "-xf", archive_path, "-C",
                                            extract_path};
      command_handler->run_command(construct);
    } else if (type == ArchiveType::ZIP) {
      // Add ZIP support for Unix-like systems
      std::vector<std::string> construct = {"unzip", "-o", archive_path, "-d",
                                            extract_path};
      command_handler->run_command(construct);
    } else {
      this->echo("Unsupported archive type");
    }

    return get_extracted_filename(name);
  }

  void git_clone(const std::string &repo_url, const std::string &clone_dir) {
    std::vector<std::string> construct = {"git", "clone", repo_url,
                                          parse_home_dir(clone_dir)};

    command_handler->run_command(construct);
  };

  void git_pull(const std::string &dir) {
    std::vector<std::string> construct = {"git", "pull"};

    this->echo("$: git pull in: " + dir);

    command_handler->run_command(construct, parse_home_dir(dir));
  };

  void create_symlinks_from_dir(const std::string source_dir,
                                const std::string &target_dir) {
    this->check_eddy_bin_path();

    std::vector<std::string> construct = {"ln -s", source_dir + "/*",
                                          target_dir};

    command_handler->run_command(construct);
    this->echo("Created symlinks for " + source_dir);
  }
}

;
#endif // SHELL_WRAPPER_H
