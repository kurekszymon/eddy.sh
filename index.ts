import { Command } from "commander";
import { languages, type Language } from "./lib/languages";
import type { Tool } from "./lib/types";

const program = new Command();

// TODO: cleanup
program
    .name('eddy.sh')
    .description('CLI to install some self container pkgs')
    .version('0.0.1');

program.command('install')
    .description('Installs a tool@version, if `version` is not specified - latest will be installed')
    .argument('<string>', 'tool')
    .argument('<string>', 'version number')
    .action(async (tool: Language, version: 'latest' | Tool['version']) => {
        const versionRegex = /(latest|^\d+\.\d+\.\d+$)/;

        if (!versionRegex.test(version)) {
            console.log('Error: Version must follow the format {number}.{number}.{number}');
            return;
        }

        if (languages[tool]) {
            const _tool = languages[tool](version);
            await _tool.install();
        }
    });

program.command('use')
    .description('Symlinks a tool@version')
    .argument('<string>', 'tool')
    .argument('<string>', 'version number')
    .action(async (tool: Language, version: 'latest' | Tool['version']) => {
        const versionRegex = /^\d+\.\d+\.\d+$/;

        if (!versionRegex.test(version)) {
            console.log('Error: Version must follow the format {number}.{number}.{number}');
            return;
        }

        if (languages[tool]) {
            const _tool = languages[tool](version);
            _tool.use();
        }
    });

program.command('delete')
    .description('deletes a tool@version')
    .argument('<string>', 'tool')
    .argument('<string>', 'version number')
    .action(async (tool: Language, version: 'latest' | Tool['version']) => {
        const versionRegex = /^\d+\.\d+\.\d+$/;

        if (!versionRegex.test(version)) {
            console.log('Error: Version must follow the format {number}.{number}.{number}');
            return;
        }

        if (languages[tool]) {
            const _tool = languages[tool](version);
            await _tool.delete();
        }
    });

program.parse();