#ifndef CONFIG_H
#define CONFIG_H

#include <iostream>
#include <memory>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <variant>
#include <vector>

#include "yaml-cpp/yaml.h"

#include "../Languages/LanguageFactory.hpp"
#include "../Shell/Shell.hpp"
#include "../Yaml/CustomScript.hpp"
#include "../Yaml/Repository.hpp"

enum ConfigItem {
  REPOSITORIES,
  CUSTOM_SCRIPTS,
};

class Config {
public:
  explicit Config(std::shared_ptr<Shell> shell) : shell(shell) {
    load_yaml_config("config.yaml");
  };

  Repositories repositories;
  std::vector<CustomScript> custom_scripts;
  std::vector<std::shared_ptr<Language>> languages;

private:
  void load_yaml_config(const std::string &yaml_file);
  std::shared_ptr<Shell> shell;
};

#endif // CONFIG_H
