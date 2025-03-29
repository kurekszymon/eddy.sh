#include <string>
#include <vector>

struct Repository {
  std::string name;
  std::string url;
};

struct Repositories {
  std::string clone_path;
  std::vector<Repository> vector;
};
