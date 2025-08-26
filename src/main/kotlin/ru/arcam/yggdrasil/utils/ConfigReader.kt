package ru.arcam.yggdrasil.utils

import java.io.File

class ConfigReader {
    companion object {
        fun loadConfig(name: String): File? {
            // Пытаемся найти конфигурационный файл в разных местах
            val configdir = File(".").listFiles()!!.firstOrNull{ it.isDirectory && it.name == "config" }
            if (configdir == null || configdir.listFiles() == null) {
                println("not found config dir")
                println(File(".").absolutePath)
                return null
            }

            val configFile = configdir.listFiles()!!.firstOrNull{it.name == name}
            return configFile
        }

    }
}