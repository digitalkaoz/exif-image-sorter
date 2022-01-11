import meow from "meow";
import chalk from 'chalk';
import { dirname, resolve } from 'path';
import { existsSync, statSync, mkdirSync } from 'fs';

export const definition = meow(chalk`
{whiteBright Usage}
  $ image-sort {yellow <src-folder>} {yellow <dest-folder>}

{whiteBright Examples}
  $ image-sort {yellow src/} {yellow dest/}
`, {
  autoHelp: true
});

export const getInputs = (inputs: string[]) => {
  let [srcDir, targetDir] = inputs.map(dir => resolve(dir));

  if (!srcDir || !targetDir) {
    definition.showHelp(1);
  }

  if (!existsSync(srcDir)) {
    console.log(chalk.red(`"${srcDir}" doesnt exists`));
    definition.showHelp(1);
  }

  if (statSync(srcDir).isFile()) {
    srcDir = dirname(srcDir);
  }

  if (!existsSync(targetDir)) {
    mkdirSync(targetDir);
  }
  return [srcDir, targetDir];
}

export default definition;
