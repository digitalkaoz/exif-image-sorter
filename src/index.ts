import cli, {getInputs} from './cli';
import process from './processor';

const [srcDir, targetDir] = getInputs(cli.input);

process(srcDir, targetDir);
