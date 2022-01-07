import { Readable, Transform } from 'stream';
import readdirp, { EntryInfo } from 'readdirp';
import { read } from 'fast-exif';
import chalk from 'chalk';

const exiftool = require('node-exiftool')


export type exifCallback = (error: Error | null, data?: string | Buffer) => void;

export const missingExif = (exif: any) => !exif || !exif.exif || (!exif.exif.DateTimeOriginal && !(exif.gps && exif.gps.GPSDateStamp) && !(exif.image && exif.image.ModifyDate));

export const readFiles = (srcDir: string): Readable => {
	return readdirp(srcDir, {
		fileFilter: ['*.jpg', '*.jpeg', '*.JPG', '*.JPEG'],
		depth: 100
	})
		.pipe(parseExif)
}

export const matchFilename = (chunk: EntryInfo): string | undefined => {
	const patterns = [
		/^([\d]{8})[_|-]/,
		/^IMG-([\d]{8})[_|-]/,
		/^IMG_([\d]{8})[_|-]/,
		/^IMG_([\d]{4}_[\d]{2}_[\d]{2})[_|-]/
	]

	const p = patterns.find(p => {
		const m = chunk.basename.match(p);
		return m && m.length >= 2
	});

	if (!p) {
		return
	}

	const dateString = chunk.basename.match(p as RegExp)[1].replace('_', '');

	const year = dateString.slice(0, 4);
	const month = dateString.slice(4, 6);
	const day = dateString.slice(6, 8);

	//TODO we should write the date back to the exif data

	return `${year}-${month}-${day}`;
}

const fastExif = (chunk: EntryInfo, callback: exifCallback) => {
	return read(chunk.fullPath)
		.then((data) => {
			if (!data) {
				return Promise.resolve({ stat: chunk, exif: null });
			}
			return Promise.resolve({
				stat: chunk,
				exif: data
			});
		})
		.catch(() => Promise.resolve({ stat: chunk, exif: null }))
}

const exifTool = (chunk: EntryInfo, callback: exifCallback) => {
	const ep = new exiftool.ExiftoolProcess()

	return ep
		.open()
		.then(() => ep.readMetadata(chunk.fullPath, ['-File:all']))
		.then((data: any, error: Error) => {
			if (!data || error || !data.data) {
				return Promise.resolve({ stat: chunk, exif: null });
			}
			ep.close();
			return Promise.resolve({
				stat: chunk,
				exif: {exif: data.data[0]}
			});
		})
		.catch(() => Promise.resolve({ stat: chunk, exif: null }))
}

export const parseExif = new Transform({
	writableObjectMode: true,
	transform: (chunk: EntryInfo, encoding: string, callback: exifCallback) => {
		// try fast-exif first
		fastExif(chunk, callback)
		.then((data) => {
			if (data && data.exif && !missingExif(data.exif)) {
				return Promise.resolve(data);
			}
			return exifTool(chunk, callback);
		})
		.then((data) => {
			return callback(null, JSON.stringify(data));
		})
		.catch((e) => {
			console.error(chalk.red(e));
			return callback(null, JSON.stringify({
				stat: chunk,
				exif: null
			}));
		});
	}
});
