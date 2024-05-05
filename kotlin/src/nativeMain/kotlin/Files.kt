import kotlinx.coroutines.Job
import kotlinx.coroutines.channels.Channel
import kotlinx.coroutines.channels.toList
import kotlinx.coroutines.coroutineScope
import kotlinx.coroutines.joinAll
import kotlinx.coroutines.launch
import kotlinx.datetime.*
import kotlinx.io.files.Path as KPath
import okio.FileSystem
import okio.IOException
import okio.Path

import okio.Path.Companion.DIRECTORY_SEPARATOR
import okio.Path.Companion.toPath
import platform.Metal.MTLCaptureDestination

val fs = FileSystem.SYSTEM

data class FileDate(val file: Path, val date: LocalDate?) {
    fun dateAsPath(destination: String): Path {
        return "${destination}$DIRECTORY_SEPARATOR${this.date?.year}$DIRECTORY_SEPARATOR${
            this.date?.month?.number.toString().padStart(2, '0')
        }$DIRECTORY_SEPARATOR${this.file.name}".toPath()
    }
}


fun readableDir(path: String): Path {
    val dir = path.toPath()

    if (!fs.exists(dir)) {
        throw Error("unreadable SRC directory.")
    }

    return dir
}

fun readDirectory(path: String): Sequence<Path> {
    return fs.listRecursively(readableDir(path)).filter { isMediaFile(it) }
}

fun isMediaFile(path: Path): Boolean {
    val meta = fs.metadata(path)
    // filter directories
    if (meta.isDirectory) {
        return false
    }
    // filter files without extension
    val ext = path.name.substringAfterLast('.', "").lowercase()
    if (ext.isEmpty()) {
        return false
    }

    //look for the extensions of interest
    val extensionsOfInterest = listOf("jpg", "mov", "mp4", "m4v", "jpeg")

    return ext in extensionsOfInterest
}

fun convertPath(path: Path): KPath {
    return KPath(path = DIRECTORY_SEPARATOR + path.segments.joinToString(DIRECTORY_SEPARATOR))
}

fun moveFile(src: Path, target: Path): Path {
    fs.createDirectories(target.parent!!)
    fs.atomicMove(src, target)

    return target
}

suspend fun resolveFileDates(src: String): List<FileDate> {
    val files = readDirectory(src)
    val resolved = Channel<FileDate>(files.count())

    coroutineScope {
        for (file in files) {
            launch {
                resolved.send(FileDate(file = file, date = readDateFromFile(file)))
            }
        }
        resolved.receive()
        resolved.close()
    }

    return resolved.toList()
}

suspend fun moveFiles(files: List<FileDate>, target: String): List<FileDate> {
    val jobs: MutableList<Job> = mutableListOf()
    coroutineScope {
        for (file in files) {
            if (file.date == null) {
                println("skipping \"${file.file}\" missing date.")
                continue
            }

            jobs.add(launch {
                try {
                    val dest = moveFile(file.file, file.dateAsPath(target))
                    println("moved ${file.file} (${file.date}) to $dest")
                } catch (e: IOException) {
                    println("could not move ${file.file}: ${e.message}")
                }
            })
        }
    }
    jobs.joinAll()

    return files
}