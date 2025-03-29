
#ifndef SHELL_WRAPPER_H
#define SHELL_WRAPPER_H

#include <boost/asio.hpp>
#include <filesystem>
#include <iostream>
#include <memory>
#include <vector>

#include "boost/filesystem.hpp"
#include "boost/process.hpp"
#include <boost/algorithm/string/join.hpp>

#include "../CommandHandler/CommandHandler.hpp"
#include "../Globals/Globals.hpp"

// #include "../Shell/ArchiveExtractor.hpp"
#include "../Yaml/CustomScript.hpp"

namespace bp = boost::process;
namespace bfs = boost::filesystem;
namespace fs = std::filesystem;

namespace {
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

std::string check_eddy_path() {
  std::string eddy_path = parse_home_dir(EDDY_PATH);

  if (!fs::exists(eddy_path)) {
    fs::create_directory(eddy_path);
  }

  return parse_home_dir(eddy_path);
}

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

  explicit ShellWrapper(std::shared_ptr<CommandHandler> handler)
      : command_handler(handler) {}

  void echo(const std::string &message) {
    command_handler->run_command("echo " + message);
  }

  void run_custom_command(CustomScript command) {
    this->echo("$ " + command.name + " started..");
    command_handler->run_command(command.cmd);
  }

  std::pair<std::string, std::string>
  curl(const std::string url, const std::string name,
       curl_type type = curl_type::download) {

    const std::string eddy_path = check_eddy_path();
    const std::string curled_filename = get_curl_filename(url);
    const std::string output_path = eddy_path + "/" + curled_filename;

    const std::vector<std::string> args = {
        "curl", "-L", "--output-dir", eddy_path, "-O", url,
    };

    std::string command = boost::algorithm::join(args, " ");

    command_handler->run_command(command);

    return std::make_pair(eddy_path, curled_filename);
  }

  void make_executable(const std::string &file_dir,
                       const std::string &filename) {
    const std::string filepath = file_dir + "/" + filename;
    this->echo("$ chmod +x " + filepath);

    auto perms =
        fs::perms::owner_exec | fs::perms::group_exec | fs::perms::others_exec;

    fs::permissions(filepath, perms, fs::perm_options::add);

    this->echo("$ File made executable: " + filepath);
  };

  void run_script_file(const std::string &file_dir,
                       const std::string &filename) {
    std::vector<std::string> args = {"sh", "-c", "./" + filename};
    std::string command = boost::algorithm::join(args, " ");

    this->command_handler->run_command(command, file_dir);
  }

  std::string tar(const std::string &path, const std::string name,
                  ArchiveType type) {
    std::string extract_path = check_eddy_path();
    std::string archive_path = path + "/" + name;
// Determine extraction command based on OS and archive type
#ifdef _WIN32
    if (type == ArchiveType::ZIP) {
      std::vector<std::string> args = {
          "powershell", "Expand-Archive",   "-Path",
          archive_path, "-DestinationPath", extract_path};
      std::string command = boost::algorithm::join(args, " ");

      this->command_handler->run_command(command);
    } else {
      this->echo("Only ZIP files supported on Windows");
    }
#else
    std::vector<std::string> args = {"tar", "-xzf", archive_path, "-C",
                                     extract_path};

    // TAR and TAR.GZ extraction
    if (type != ArchiveType::TAR_GZ || type != ArchiveType::TAR) {
      this->echo("Unsupported archive type on Unix-like system");
    }
    std::string command = boost::algorithm::join(args, " ");

    this->command_handler->run_command(command);
#endif

    return get_extracted_filename(name);
  }

  void bootstrap_cmake(const std::string &cmake_source_path) {
    const std::string cmake_path = check_eddy_path() + "/" + cmake_source_path;
    const std::vector<std::string> args = {"sh", "-c", "./bootstrap"};
    std::string command = boost::algorithm::join(args, " ");
    this->command_handler->run_command(command, cmake_path);
  }

  void git_clone(const std::string &repo_url, const std::string &clone_dir) {
    std::vector<std::string> args = {"git", "clone", repo_url};
    std::string command = boost::algorithm::join(args, " ");

    this->echo("$: Starting cloning " + repo_url + " to " + clone_dir);

    this->command_handler->run_command(command, parse_home_dir(clone_dir));
  };
}

//   std::string tar(const std::string &path, ArchiveType type) {
//     const std::string eddy_path = check_eddy_path();
//     // std::string unpacked_filename = extractor.unpack(path,
//     eddy_path, type);

//     // return unpacked_filename;
//     return "";
//   }

//   void bootstrap_cmake(const std::string &cmake_source_path) {
//     const std::string cmake_path = check_eddy_path() + "/" +
//     cmake_source_path;

//     bp::ipstream std_out;
//     bp::ipstream std_err;
//     bp::child c("sh -c \"echo 12323 && pwd && echo hello\"",
//                 bp::start_dir(cmake_path), bp::std_out > std_out,
//                 bp::std_err > std_err);
//     c.wait();

//     // pipe_output(std_out, logger, "bootstrap_cmake:");
//     // pipe_output(std_err, logger, "bootstrap_cmake:");
//   }

//   void echo(const std::string msg) { logger->update(msg); }
// }
;
#endif // SHELL_WRAPPER_H
