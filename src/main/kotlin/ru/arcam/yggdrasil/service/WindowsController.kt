package ru.arcam.yggdrasil.service
import ru.arcam.yggdrasil.leaf.Leaf
import java.io.BufferedReader
import java.io.InputStreamReader
import java.lang.Exception

class WindowsController(leaf: Leaf) : IController(leaf) {
    override fun stop(): String {
        val SERVER_STATUS = arrayOf("sc", "stop", leaf.name)
        try {
            val process = ProcessBuilder(*SERVER_STATUS).redirectErrorStream(true).start()
            process.waitFor()
            return "OK"
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "ERROR"
    }
    
    override fun start(): String {
        val SERVER_STATUS = arrayOf("sc", "start", leaf.name)
        try {
            val process = ProcessBuilder(*SERVER_STATUS).redirectErrorStream(true).start()
            process.waitFor()
            return "OK"
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "ERROR"
    }

    override fun restart(): String {
        val stopStatus = stop()
        if (stopStatus == "OK") {
            val startStatus = start()
            if (startStatus == "OK")
                return "OK"
        }
        return "ERROR"
    }
    
    override fun status(): String {
        val SERVER_STATUS = arrayOf("sc", "query", leaf.name)
        var status = "UNAVAIVABLE"
        try {
            val process = ProcessBuilder(*SERVER_STATUS).redirectErrorStream(true).start()
            process.waitFor()
            val input = process.inputStream
            val br = BufferedReader(InputStreamReader(input))
            var cnt = 0
            br.forEachLine {
                if (!it.isEmpty()) {
                    cnt++
                    if (cnt == 3)
                        status = it.split(':').last().replace(Regex(" \\d* "), "").trim()
                }
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return status
    }

    override fun callMethod(method: String, args: List<String>): String {
        return "OK"
    }
}
