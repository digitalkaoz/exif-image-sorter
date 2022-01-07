import { EntryInfo } from 'readdirp';
import { Writable } from 'stream';
import { promises as fsPromises } from 'fs';
import chalk from 'chalk';
import { matchFilename, missingExif } from './reader';
import { resolve } from 'path';

const piexif = require("piexifjs");

class Writer extends Writable {
	constructor(private dir: string) {
		super();
	}

	private dateFromExif(exif): Date {
		if (exif.exif && exif.exif.DateTimeOriginal) {
			return this.dateFromString(exif.exif.DateTimeOriginal);
		} else if (exif.gps && exif.gps.GPSDateStamp) {
			return this.dateFromString(exif.gps.GPSDateStamp.replace(':', '-'));
		} else if (exif.image && exif.image.ModifyDate) {
			return this.dateFromString(exif.image.ModifyDate);
		}
		console.error(exif);
		throw Error("unable to get date from exifdata")

	}
	private dateFromString(dateString: string): Date {

		let date = new Date(dateString);
		if (isNaN(date as any)) {
			date = new Date(dateString.split(" ")[0].replace(':','-'));
		}
		if (isNaN(date as any)) {
			console.error(chalk.red(`unable to parse date string "${dateString}"`))
			return null;
		}

		return date;
	}

	private parseDir(dateString: string): string {
		const date = this.dateFromString(dateString);

		const year = date.getFullYear();
		const month = (date.getMonth() + 1).toString().padStart(2, '0');

		return `${this.dir}/${year}/${month}`;
	}

	private createDir(dir: string): Promise<any> {
		return fsPromises.mkdir(dir, {
			recursive: true
		});
	}

	private copyFile(stat: EntryInfo, exif: any, callback: (error?: Error | null) => void) {
		if (!exif) {
			return callback(null);
		}
		const targetDir = this.dirFromExif(exif);

		if (!targetDir) {
			return callback(null);
		}

		if (resolve(targetDir) === resolve(stat.fullPath.replace(`/${stat.basename}`, ''))) {
			console.log(chalk.dim(`${chalk.yellow(stat.fullPath)} already in place`));
			const d = this.dateFromExif(exif);
			if (d) {
				return fsPromises.utimes(stat.fullPath, d, d).then(() => callback(null));
			} else {
				return callback(null)
			}
		}

		this.createDir(targetDir as string)
			.then(() => fsPromises.copyFile(stat.fullPath, `${targetDir}/${stat.basename}`))
			.then(() => fsPromises.unlink(stat.fullPath))
			.then(() => {
				console.log(chalk`moved {yellow ${stat.fullPath}} to {green ${targetDir as string}/${stat.basename}}`);
				const d = this.dateFromExif(exif);
				if (d) {
					fsPromises.utimes(`${targetDir}/${stat.basename}`, d, d).then(() => callback(null));
				} else {
					callback(null)
				}
			})
			.catch((e) => {
				console.log(e);
				//callback(null)
			})
	}

	private async addMissingMetadata(stat: EntryInfo, date: Date) {
		const exif = {};
		exif[piexif.ExifIFD.DateTimeOriginal] = date.toISOString();
		const exifBytes = piexif.dump({"Exif":exif});

		return fsPromises.readFile(stat.fullPath)
			.then(data => Promise.resolve(data.toString("binary")))
			.then(data => Promise.resolve(piexif.insert(exifBytes, data)))
			.then(data => Promise.resolve(Buffer.from(data, "binary")))
			.then(data => fsPromises.writeFile(stat.fullPath, data))
			.then(() => console.log(chalk`fixed metadata for {yellow ${stat.fullPath}}`))
			.then(() => ({
				exif: {
					DateTimeOriginal : date.toISOString()
				}
			}))
			.catch((e) => {
				console.error(e);
				return Promise.resolve()
			})
	}

	private dirFromExif(exif: any) {
		if (exif.exif && exif.exif.DateTimeOriginal) {
			return this.parseDir(exif.exif.DateTimeOriginal);
		} else if (exif.gps && exif.gps.GPSDateStamp) {
			return this.parseDir(exif.gps.GPSDateStamp.replace(':', '-'));
		} else if (exif.image && exif.image.ModifyDate) {
			return this.parseDir(exif.image.ModifyDate);
		}
		console.error(exif);
		throw Error("unable to get date from exifdata")
	}

	private handleMissingExif(stat: EntryInfo, exif: any, callback: (error?: Error | null) => void) {
		const date = matchFilename(stat);
		if (date) {
			return this.addMissingMetadata(stat, new Date(date)).then((exif) => {
				this.copyFile(stat, exif, callback);
			});
		}

		console.log(chalk.dim(`cant extract metadata from ${chalk.yellow(stat.fullPath)}`));
		return callback(null);
	}

	_write(data: Buffer, encoding: string, callback: (error?: Error | null) => void) {
		let { stat, exif }: { stat: EntryInfo, exif: any } = JSON.parse(data.toString());

		if (missingExif(exif)) {
			return this.handleMissingExif(stat, exif, callback);
		}

		this.copyFile(stat, exif, callback);
	}
}

export const writeFiles = (targetDir: string) => new Writer(targetDir);

