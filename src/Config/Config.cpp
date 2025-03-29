#include "Config.hpp"

void Config::load_yaml_config(const std::string &yaml_file) {
  YAML::Node config = YAML::LoadFile(yaml_file);
  std::unique_ptr factory = std::make_unique<LanguageFactory>(shell);
  // handle empty yaml?

  if (config["repositories"]) {
    auto repos_node = config["repositories"];

    if (repos_node["clone_path"]) {
      repositories.clone_path = repos_node["clone_path"].as<std::string>();
    }

    for (const auto &yaml_repo : repos_node["repos"]) {
      auto node = yaml_repo.begin();

      Repository repo;
      repo.name = node->first.as<std::string>();
      repo.url = node->second.as<std::string>();
      repositories.vector.push_back(repo);
    }
  }
  if (config["custom_scripts"]) {
    auto custom_scripts_node = config["custom_scripts"];

    for (const auto &yaml_script : custom_scripts_node) {
      auto node = yaml_script.begin();

      CustomScript script;
      script.name = node->first.as<std::string>();
      script.cmd = node->second.as<std::string>();
      custom_scripts.push_back(script);
    }
  }
  if (config["languages"]) {
    auto languages_node = config["languages"];

    for (const auto &yaml_language : languages_node) {
      auto node = yaml_language.begin();

      std::string language_name = node->first.as<std::string>();
      std::shared_ptr<Language> lang = factory->create(language_name);

      auto tools = node->second;

      if (tools.size() == 0) {
        continue;
      }

      for (const auto &yaml_tool : tools) {
        auto node = yaml_tool.begin();

        std::string name = node->first.as<std::string>();
        std::string version = node->second.as<std::string>();

        lang->load_tool(name, version);
        // move creating language factory to main maybe?
        // this way I can skip dependency injection to config
      }

      languages.push_back(lang);
    }
  }
}
