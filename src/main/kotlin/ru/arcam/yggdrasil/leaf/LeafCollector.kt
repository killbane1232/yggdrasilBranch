package ru.arcam.yggdrasil.leaf

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Service
import ru.arcam.yggdrasil.branch.BranchInfo
import ru.arcam.yggdrasil.service.*
import ru.arcam.yggdrasil.utils.ConfigReader
import ru.arcam.yggdrasil.utils.NameResolver
import ru.arcam.yggdrasil.ws.TrunkConnection
import java.io.File
import java.util.*
import java.util.stream.Collectors
import kotlin.collections.ArrayList
import kotlin.collections.HashMap
import kotlin.reflect.typeOf


@Service
class LeafCollector {
    val isWindows: Boolean = System.getProperty("os.name").lowercase(Locale.getDefault()).contains("win")
    private var leafStatus: Map<String, String> = HashMap()
    var lock: Any = Any()
    var linkedServices: ArrayList<Leaf> = ArrayList()
    var configuredServices: ArrayList<Leaf> = ArrayList()
    val serviceName = NameResolver.name

    private fun findConfigFile(filename: String): File? {
        val possiblePaths = listOf(
            File("config/$filename"),  // Относительно текущей директории
            File("./config/$filename"), // Явно относительно текущей директории
            File("/app/config/$filename"), // Абсолютный путь от текущей директории
        )
        return possiblePaths.firstOrNull { it.exists() }
    }

    @Scheduled(fixedRate = 5000)
    fun updateConfiguredServices() {
        try {
            var configFile = ConfigReader.loadConfig("leaves.config")
            var leafServices: List<String> = ArrayList()
            var confServicesBuffer = ArrayList<Leaf>()
            val newServices = ArrayList<Leaf>()

            synchronized(lock) {
                confServicesBuffer = ArrayList(configuredServices.stream()
                    .filter({x -> x.controller is CustomController})
                    .collect(Collectors.toList()))
            }

            if (configFile != null) {
                leafServices = configFile.readLines()
                    .filter { it.isNotBlank() && !it.startsWith("#") }
            }
            leafServices.forEach { x ->
                val newLeaf = Leaf(x, "UNAVAIVABLE", serviceName, ArrayList(), null)
                newLeaf.controller = if (isWindows) WindowsController(newLeaf) else LinuxController(newLeaf)
                newLeaf.status = newLeaf.controller!!.status()
                confServicesBuffer.add(newLeaf)
            }
            configFile = findConfigFile("docker.config")
            leafServices = ArrayList()

            if (configFile != null) {
                leafServices = configFile.readLines()
                    .filter { it.isNotBlank() && !it.startsWith("#") }
            }
            leafServices.forEach { x ->
                val newLeaf = Leaf(x, "UNAVAIVABLE", serviceName, ArrayList(), null)
                newLeaf.controller = DockerController(newLeaf)
                newLeaf.status = newLeaf.controller!!.status()
                confServicesBuffer.add(newLeaf)
            }

            synchronized(lock) {
                configuredServices = confServicesBuffer
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
    }

    @Scheduled(fixedRate = 5000)
    fun reportLeaves() {
        var confServices: ArrayList<Leaf>
        synchronized(lock) {
            confServices = ArrayList(configuredServices.union(linkedServices))
        }

        val currentInfo = BranchInfo(serviceName, confServices)

        TrunkConnection.WSClient?.send("/app/assing", currentInfo)
    }

    fun callServiceMethod(serviceName: String, method: String, args: List<String>) : String {
        val leaf = configuredServices.union(linkedServices).find { x -> x.name == serviceName }
        return when(method) {
            "STATUS" -> {
                if (leaf != null) {
                    leaf.status = leaf.controller?.status()?:"UNAVAIVABLE"
                    leaf.status
                } else {
                    "UNAVAIVABLE"
                }
            }
            "START" -> leaf?.controller?.start()?:"ERROR"
            "STOP" -> leaf?.controller?.stop()?:"ERROR"
            "RESTART" -> leaf?.controller?.restart()?:"ERROR"
            "TAIL" -> leaf?.controller?.logs()?:"ERROR"
            "TAIL_N" -> leaf?.controller?.logs(args)?:"ERROR"
            else -> leaf?.controller?.callMethod(method, args)?:"ERROR"
        }
    }
}