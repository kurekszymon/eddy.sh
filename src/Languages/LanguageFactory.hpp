#ifndef LANGUAGEFACTORY_H
#define LANGUAGEFACTORY_H

#include "../ShellWrapper/ShellWrapper.hpp"
#include "Cpp.hpp"
#include "Javascript.hpp"
#include "Language.hpp"
#include "Python.hpp"

/**
 * @class LanguageFactory
 * @brief Factory class to create different language objects.
 *
 * This class is responsible for creating instances of languages like Python and
 * JavaScript.
 */
class LanguageFactory {
public:
  LanguageFactory(std::shared_ptr<ShellWrapper> shell) : shell_(shell) {
    language_factory_map["python"] = [&]() {
      return LanguageFactory::create_language<Python>();
    };
    language_factory_map["js"] = [&]() {
      return LanguageFactory::create_language<Javascript>();
    };
    language_factory_map["cpp"] = [&]() {
      return LanguageFactory::create_language<Cpp>();
    };
  }
  // Factory method to create the language instance based on the string input
  std::shared_ptr<Language> create(const std::string &language_name) {
    if (language_factory_map.find(language_name) !=
        language_factory_map.end()) {

      return language_factory_map[language_name]();
    }
    std::cout << "Language " << language_name << " is not supported!"
              << std::endl;
    return nullptr;
  }

  template <typename T> std::shared_ptr<Language> create_language() {
    return std::make_unique<T>(shell_); // Create the specific language instance
  };

private:
  std::shared_ptr<ShellWrapper> shell_;
  std::unordered_map<std::string, std::function<std::shared_ptr<Language>()>>
      language_factory_map;
};

#endif // LANGUAGEFACTORY_H