import kotlinx.coroutines.*

fun main(args: Array<String>): Unit = runBlocking {
    try {
        if (args.size != 2) {
            throw Exception("Invalid number of arguments. \n usage: \nImageSorter.kexe SRC DEST")
        }

        val resolved = resolveFileDates(args.first())
        moveFiles(resolved, args.last())
    } catch (e: Exception) {
        println(e.message)
    }
}
