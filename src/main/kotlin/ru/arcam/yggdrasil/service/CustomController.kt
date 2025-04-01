package ru.arcam.yggdrasil.service
import ru.arcam.yggdrasil.leaf.Leaf
import ru.arcam.yggdrasil.leaf.LeafCollector
import ru.arcam.yggdrasil.leaf.LeafController
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

    override fun callMethod(method: String, args: List<String>): String {
        try {
            return LeafController.leafController.callLeafMethod(leaf.name, method, args)
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return "ERROR"
    }
}
