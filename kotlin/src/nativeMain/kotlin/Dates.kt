import com.ashampoo.kim.Kim
import com.ashampoo.kim.common.ImageReadException
import com.ashampoo.kim.format.ImageMetadata
import com.ashampoo.kim.format.tiff.constant.ExifTag
import com.ashampoo.kim.format.tiff.constant.TiffTag
import com.ashampoo.kim.format.tiff.taginfo.TagInfo
import com.ashampoo.kim.kotlinx.readMetadata
import kotlinx.datetime.Instant
import kotlinx.datetime.LocalDate
import kotlinx.datetime.TimeZone
import kotlinx.datetime.format.DateTimeComponents
import kotlinx.datetime.format.FormatStringsInDatetimeFormats
import kotlinx.datetime.format.byUnicodePattern
import kotlinx.datetime.toLocalDateTime
import okio.IOException
import okio.Path

val filenamePatterns = listOf(
    "[_|-]([\\d]{8})[_|-]",                           //matches e.g. IMG_20221030-foo.jpg
    "[_|-]([\\d]{4}[_|-][\\d]{2}[_|-][\\d]{2})[_|-]", //matches e.g. IMG-2022-10-30_bar.jpg
    "^([\\d]{8})[_|-]",                               //matches e.g. 20221030-foo.jpg
    "^([\\d]{4}[_|-][\\d]{2}[_|-][\\d]{2})[_|-]",     //matches e.g. 2022-10-30_bar.jpg
)


@OptIn(FormatStringsInDatetimeFormats::class)
val exifDateFormat = DateTimeComponents.Format {
    byUnicodePattern("uuuu:MM:dd HH:mm:ss")
}
@OptIn(FormatStringsInDatetimeFormats::class)
val fileDateFormat = DateTimeComponents.Format {
    byUnicodePattern("uuuuMMdd")
}

fun readDateFromFile(file:Path): LocalDate? {
    var d = readDateFromMetaData(file)
    if (d != null) return d

    d = readDateFromFilename(file)
    if (d != null) return d

    return readDateFromFilesystem(file)
}

fun readDateFromMetaData(file: Path): LocalDate? {
    val meta: ImageMetadata?
    try {
        meta = Kim.readMetadata(convertPath(file))
    } catch (e: ImageReadException) {
        println("reading image metadata for $file failed: ${e.message}")
        return null
    }

    if (meta == null) {
        return null
    }

    val exifCreated = extractDateFromMetadata(meta, ExifTag.EXIF_TAG_DATE_TIME_ORIGINAL)
    if (exifCreated != null) {
        return exifCreated
    }

    val tiffCreated = extractDateFromMetadata(meta, TiffTag.TIFF_TAG_DATE_TIME)
    if (tiffCreated != null) {
        return tiffCreated
    }

    return null
}

fun readDateFromFilename(file: Path): LocalDate? {
    for (p in filenamePatterns) {
        val r = Regex(p)

        val match = r.find(file.name) ?: continue

        val date = match.value.replace("_","").replace("-","")

        return fileDateFormat.parse(date).toLocalDate()
    }

    return null
}

fun readDateFromFilesystem(file:Path): LocalDate? {
    try {
        val ms = fs.metadata(file).createdAtMillis ?: return null
        val instant = Instant.fromEpochMilliseconds(ms)

        return instant.toLocalDateTime(TimeZone.UTC).date
    } catch (e: IOException) {
        println("reading file metadata for $file failed: ${e.message}")
        return null
    }
}

fun extractDateFromMetadata(metadata: ImageMetadata, tag: TagInfo): LocalDate? {
    val d = metadata.findStringValue(tag)
    if (d != null) {
        return exifDateFormat.parse(d).toLocalDate()
    }
    return null
}
