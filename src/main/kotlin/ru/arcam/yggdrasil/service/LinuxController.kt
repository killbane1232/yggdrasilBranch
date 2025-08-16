package ru.arcam.yggdrasil.service
import ru.arcam.yggdrasil.leaf.Leaf
import java.io.BufferedReader
import java.io.InputStreamReader

class LinuxController(leaf: Leaf) : IController(leaf) {
    override fun stop(): String {
        val SERVER_STATUS = arrayOf("sudo", "/usr/bin/systemctl", "stop", leaf.name)
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
        val SERVER_STATUS = arrayOf("sudo", "/usr/bin/systemctl", "start", leaf.name)
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
        val SERVER_STATUS = arrayOf("sudo", "/usr/bin/systemctl", "restart", leaf.name)
        try {
            val process = ProcessBuilder(*SERVER_STATUS).redirectErrorStream(true).start()
            process.waitFor()
            return "OK"
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "ERROR"
    }
    
    override fun status(): String {
        val SERVER_STATUS = arrayOf("sudo", "/usr/bin/systemctl", "is-active", leaf.name)
        try {
            val process = ProcessBuilder(*SERVER_STATUS).redirectErrorStream(true).start()
            process.waitFor()
            val input = process.inputStream
            val br = BufferedReader(InputStreamReader(input))
            return if (br.readLine() == "active")  "RUNNING" else "STOPPED"
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "UNAVAIVABLE"
    }

    override fun logs(args: List<String>?): String {
        var num = 10
        if (args != null && args.size > 0) {
            val tmp = args[0].toIntOrNull()
            if (tmp != null && tmp > 0)
                num = tmp
        }
        val SERVER_LOGS = arrayOf("sudo", "/usr/bin/journalctl", "--unit=${leaf.name}", "-n", num.toString(), "--no-pager")
        try {
            val process = ProcessBuilder(*SERVER_LOGS).redirectErrorStream(true).start()
            process.waitFor()
            val input = process.inputStream
            val br = BufferedReader(InputStreamReader(input))
            val result = br.readLines().joinToString { it -> "\n" + it }
            return result
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "UNAVAIVABLE"
    }

    override fun callMethod(method: String, args: List<String>): String {
        return "OK"
    }
}
