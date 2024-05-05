plugins {
    kotlin("multiplatform") version "1.9.23"
}

group = "net.digitalkaoz"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

kotlin {
    val hostOs = System.getProperty("os.name")
    val isArm64 = System.getProperty("os.arch") == "aarch64"
    val isMingwX64 = hostOs.startsWith("Windows")
    val nativeTarget = when {
        hostOs == "Mac OS X" && isArm64 -> macosArm64("native")
        hostOs == "Mac OS X" && !isArm64 -> macosX64("native")
        hostOs == "Linux" && isArm64 -> linuxArm64("native")
        hostOs == "Linux" && !isArm64 -> linuxX64("native")
        isMingwX64 -> mingwX64("native")
        else -> throw GradleException("Host OS is not supported in Kotlin/Native.")
    }

    nativeTarget.apply {
        binaries {
            executable {
                entryPoint = "main"
                runTask?.run {
                    val args = providers.gradleProperty("runArgs")
                    argumentProviders.add(CommandLineArgumentProvider {
                        args.orNull?.split(' ') ?: emptyList()
                    })
                }
            }
        }
    }
    sourceSets {
        val nativeMain by getting
        val nativeTest by getting

        val ioVersion = "0.3.3"
        val kimVersion = "0.17.7"
        val okioVersion = "3.9.0"
        val dateVersion = "0.6.0-RC.2"
        val coroutinesVersion = "1.8.1-Beta"


        val commonMain by getting {
            dependencies {
                implementation("org.jetbrains.kotlinx:kotlinx-io-core:$ioVersion")
                implementation("com.squareup.okio:okio:$okioVersion")
                implementation("com.ashampoo:kim:$kimVersion")
                implementation("org.jetbrains.kotlinx:kotlinx-datetime:$dateVersion")
                implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:$coroutinesVersion")
            }
        }
        val commonTest by getting {
            dependencies {
                implementation("com.squareup.okio:okio-fakefilesystem:$okioVersion")
            }
        }
    }
}
