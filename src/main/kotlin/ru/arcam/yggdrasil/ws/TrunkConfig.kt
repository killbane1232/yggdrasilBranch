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

    private fun loadConfig() {
        try {
            val configFile = File("websocket.config")

            if (configFile.exists()) {
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
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
    }

    fun getWebSocketUrl(): String = "ws://$host:$port$path"
    fun getTimeout(): Long = timeout
} 