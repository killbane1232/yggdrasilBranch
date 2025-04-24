package ru.arcam.yggdrasil.ws

import java.io.File
import java.net.URI

class TrunkConfig {
    private var host: String = "localhost"
    private var port: Int = 8080
    private var path: String = "/ws"
    private var timeout: Long = 5000

    init {
        loadConfig()
    }

    fun loadConfig() {
        try {
            // Пытаемся найти конфигурационный файл в разных местах
            val configdir = File(".").listFiles()!!.firstOrNull{ it.isDirectory && it.name == "config" }
            if (configdir == null || configdir.listFiles() == null) {
                println("not found config dir")
                return
            }

            val configFile = configdir.listFiles()!!.firstOrNull{it.name == "websocket.config"}

            if (configFile != null) {
                println("found config: " + configFile.absolutePath)
                configFile.readLines().forEach { line ->
                    if (line.isNotBlank() && !line.startsWith("#")) {
                        val parts = line.split("=")
                        if (parts.size == 2) {
                            when (parts[0].trim()) {
                                "websocket.host" -> host = parts[1].trim()
                                "websocket.port" -> port = parts[1].trim().toInt()
                                "websocket.path" -> path = parts[1].trim()
                                "websocket.timeout" -> timeout = parts[1].trim().toLong()
                            }
                        }
                    }
                }
                println("config after reading: " + getWebSocketUrl())
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
    }

    fun getWebSocketUrl(): String = "ws://$host:$port$path"
    fun getTimeout(): Long = timeout
} 