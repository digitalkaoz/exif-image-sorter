import {pipeline} from 'stream';
import chalk from 'chalk';

import {readFiles} from './reader';
import {writeFiles } from './writer';

const errorHandler = (e:Error) => console.error(chalk.red(e.message));

export default (srcDir: string, targetDir: string) => {
	pipeline(
		readFiles(srcDir).on('error', errorHandler),
		writeFiles(targetDir).on('error', errorHandler),
		(error) => {
			if (error) {
				console.error(chalk.red(JSON.stringify(error)));
			} else {
				console.log(chalk.yellow`done`);
			}
		}
	)
}
