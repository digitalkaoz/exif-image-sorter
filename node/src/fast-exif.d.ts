declare module 'fast-exif' {
	import { Readable } from 'stream';

	function read(
		file: String,
		fullRead?: boolean
	  ): Promise<Object>;

}

declare module 'exif-parser' {
	import { Readable } from 'stream';

	class Parser {
		parse(): object
	}

	function create(
		buffer: Readable,
	  ): Parser;
}
