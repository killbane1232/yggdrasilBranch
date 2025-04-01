package ru.arcam.yggdrasil.service
import ru.arcam.yggdrasil.leaf.Leaf
import java.io.BufferedReader
import java.io.InputStreamReader

class DockerController(leaf: Leaf) : IController(leaf) {
    override fun stop(): String {
        val SERVER_STATUS = arrayOf("docker", "stop", leaf.name)
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
        val SERVER_STATUS = arrayOf("docker", "start", leaf.name)
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
        val SERVER_STATUS = arrayOf("docker", "restart", leaf.name)
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
        val SERVER_STATUS = arrayOf("docker", "ps", "-f", "name=" + leaf.name)
        try {
            val process = ProcessBuilder(*SERVER_STATUS).redirectErrorStream(true).start()
            process.waitFor()
            val input = process.inputStream
            val br = BufferedReader(InputStreamReader(input))
            return if (br.readLine() != null)  "RUNNING" else "STOPPED"
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "UNAVAIVABLE"
    }

    override fun callMethod(method: String, args: List<String>): String {
        return "OK"
    }
}
