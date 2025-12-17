import { Command } from "commander";
import { cpp } from "./lib/languages/cpp";
import type { Tool } from "./lib/types";

const program = new Command();

program
    .name('eddy.sh')
    .description('CLI to install some self container pkgs')
    .version('0.0.1');

program.command('install')
    .description('Installs a tool@version, if `version` is not specified - latest will be installed')
    .argument('<string>', 'tool')
    .option('-v, --version <string>', 'version number', 'latest')
    .action((tool: keyof typeof cpp, options: { version: 'latest' | Tool['version']; }) => {
        const versionRegex = /(latest|^\d+\.\d+\.\d+$)/;

        if (!versionRegex.test(options.version)) {
            console.log('Error: Version must follow the format {number}.{number}.{number}');
            return;
        }

        if (cpp[tool]) {
            const _tool = cpp[tool](options.version);
            _tool.install();
        }
    });

program.command('use')
    .description('Symlinks a tool@version')
    .argument('<string>', 'tool')
    .requiredOption('-v, --version <string>', 'version number')
    .action((tool: keyof typeof cpp, options: { version: 'latest' | Tool['version']; }) => {
        const versionRegex = /^\d+\.\d+\.\d+$/;

        if (!versionRegex.test(options.version)) {
            console.log('Error: Version must follow the format {number}.{number}.{number}');
            return;
        }

        if (cpp[tool]) {
            const _tool = cpp[tool](options.version);
            _tool.use();
        }
    });

program.parse();