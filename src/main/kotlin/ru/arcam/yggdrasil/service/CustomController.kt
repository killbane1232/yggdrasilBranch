package ru.arcam.yggdrasil.service
import ru.arcam.yggdrasil.leaf.*
import java.io.BufferedReader
import java.io.InputStreamReader
import java.lang.Exception

class CustomController(leaf: Leaf, val leafCollector: LeafCollector) : IController(leaf) {
    override fun stop(): String {
        callMethod("stop", listOf())
        return "OK"
    }
    
    override fun start(): String {
        callMethod("start", listOf())
        return "OK"
    }

    override fun restart(): String {
        callMethod("restart", listOf())
        return "OK"
    }
    
    override fun status(): String {
        var status = "UNAVAIVABLE"
        try {
            status = callMethod("status", listOf())
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return status
    }

    override fun logs(args: List<String>?): String {
        return "UNAVAIVABLE"
    }

    override fun callMethod(method: String, args: List<String>): String {
        println("CALL ME LATER $method ${leaf.name} ${args.joinToString { it }}")
        try {
            val request = LeafMethodMessage(method, args)
            return LeafHttpController.leafHttpController.callLeafMethod(leaf, request).body!!
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "ERROR"
    }
}
