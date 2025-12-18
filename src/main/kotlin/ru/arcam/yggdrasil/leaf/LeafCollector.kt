package ru.arcam.yggdrasil.leaf

import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.module.kotlin.readValue
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Service
import ru.arcam.yggdrasil.branch.BranchInfo
import ru.arcam.yggdrasil.service.*
import ru.arcam.yggdrasil.user.UserRight
import ru.arcam.yggdrasil.utils.ConfigReader
import ru.arcam.yggdrasil.utils.NameResolver
import ru.arcam.yggdrasil.ws.TrunkConnection
import java.io.File
import java.util.*
import java.util.stream.Collectors
import kotlin.collections.ArrayList
import kotlin.collections.HashMap


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
                val strSplt = x.split(';')
                val allowedUsers = HashMap<String, UserRight>()
                if (strSplt.size > 1) {
                    val userInfo = strSplt[1].split(',')
                    for (it in userInfo) {
                        val itemInfo = it.split(':')
                        val userName = itemInfo[0]
                        val rights = UserRight.getFromString(itemInfo[1])
                        allowedUsers[userName] = rights
                    }
                }
                val newLeaf = Leaf(strSplt[0], "UNAVAIVABLE", serviceName, ArrayList(), allowedUsers, null)
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
                val strSplt = x.split(';')
                val allowedUsers = HashMap<String, UserRight>()
                if (strSplt.size > 1) {
                    val userInfo = strSplt[1].split(',')
                    for (it in userInfo) {
                        val itemInfo = it.split(':')
                        val userName = itemInfo[0]
                        val rights = UserRight.getFromString(itemInfo[1])
                        allowedUsers[userName] = rights
                    }
                }
                val newLeaf = Leaf(strSplt[0], "UNAVAIVABLE", serviceName, ArrayList(), allowedUsers, null)
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
        var userConfigFile = ConfigReader.loadConfig("user.config")
        var users: List<String> = listOf()
        if (userConfigFile != null) {
            users = userConfigFile.readLines()
                .filter { it.isNotBlank() && !it.startsWith("#") }
        }
        val allowedUsers = HashMap<String, UserRight>()
        users.forEach { x ->
            val data = x.split(':')
            if (data.size > 1) {
                allowedUsers[data[0]] = UserRight.getFromString(data[1])
            }
        }

        val currentInfo = BranchInfo(serviceName, confServices, allowedUsers)

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
            "RIGHTS" -> handleRightsCommand(serviceName, args)
            else -> leaf?.controller?.callMethod(method, args)?:"ERROR"
        }
    }

    private fun handleRightsCommand(leafName: String, args: List<String>): String {
        return try {
            if (args.isEmpty()) {
                return "ERROR: No JSON provided"
            }

            // Объединяем все аргументы обратно в JSON строку
            // args содержит части JSON, разделённые по ':'
            val jsonString = args.joinToString(":")
            println("Got message $jsonString")
            val objectMapper = ObjectMapper()
            val rightsMap: Map<String, UserRight> = objectMapper.readValue(jsonString)

            synchronized(lock) {
                if (leafName == "NULL") {
                    // Устанавливаем глобальные права
                    saveGlobalRights(rightsMap)
                    "OK: Global rights updated"
                } else {
                    // Устанавливаем права для конкретного листа
                    val leaf = configuredServices.union(linkedServices).find { x -> x.name == leafName }
                    if (leaf != null) {
                        updateLeafRights(leaf, rightsMap)
                        "OK: Rights updated for leaf $leafName"
                    } else {
                        "ERROR: Leaf $leafName not found"
                    }
                }
            }
        } catch (e: Exception) {
            e.printStackTrace()
            "ERROR: ${e.message}"
        }
    }

    private fun saveGlobalRights(newRightsMap: Map<String, UserRight>) {
        // Используем тот же путь, что и ConfigReader
        val configDir = File("config")
        if (!configDir.exists()) {
            configDir.mkdirs()
        }
        
        val userConfigFile = File(configDir, "user.config")
        
        // Читаем существующие права из файла
        val existingRights = HashMap<String, UserRight>()
        if (userConfigFile.exists()) {
            val users = userConfigFile.readLines()
                .filter { it.isNotBlank() && !it.startsWith("#") }
            users.forEach { x ->
                val data = x.split(':')
                if (data.size > 1) {
                    existingRights[data[0]] = UserRight.getFromString(data[1])
                }
            }
        }
        
        // Объединяем существующие права с новыми (новые перезаписывают старые)
        val mergedRights = existingRights.toMutableMap()
        mergedRights.putAll(newRightsMap)
        
        // Сохраняем объединённые права
        val lines = mergedRights.map { (userName, rights) ->
            val rightsString = buildString {
                if (rights.read) append('r')
                if (rights.write) append('w')
                if (rights.execute) append('x')
            }
            "$userName:$rightsString"
        }
        println("updated ${lines.joinToString { "," }}")
        userConfigFile.writeText(lines.joinToString("\n"))
    }

    private fun updateLeafRights(leaf: Leaf, newRightsMap: Map<String, UserRight>) {
        synchronized(lock) {
            // Определяем, в каком конфигурационном файле находится лист
            val configFileName = when (leaf.controller) {
                is DockerController -> "docker.config"
                is LinuxController, is WindowsController -> "leaves.config"
                else -> null // CustomController - лист из linkedServices, не сохраняем в файл
            }
            
            // Получаем существующие права листа
            val existingRights = leaf.allowedUsers.toMutableMap()
            
            // Если лист из конфигурационного файла, читаем права из файла для синхронизации
            if (configFileName != null) {
                val fileRights = readLeafRightsFromConfigFile(leaf.name, configFileName)
                if (fileRights != null) {
                    // Используем права из файла как основу (они более актуальны)
                    existingRights.clear()
                    existingRights.putAll(fileRights)
                }
            }
            
            // Объединяем существующие права с новыми (новые перезаписывают старые)
            existingRights.putAll(newRightsMap)
            
            // Обновляем конфигурационный файл, если лист из configuredServices
            if (configFileName != null) {
                updateLeafConfigFile(leaf.name, existingRights, configFileName)
            }
            
            // Создаём новый объект Leaf с обновлёнными правами
            val updatedLeaf = Leaf(
                name = leaf.name,
                status = leaf.status,
                attachedBranch = leaf.attachedBranch,
                hooks = leaf.hooks,
                allowedUsers = existingRights,
                controller = leaf.controller,
                port = leaf.port,
                url = leaf.url
            )

            // Обновляем лист в соответствующей коллекции
            val indexInConfigured = configuredServices.indexOfFirst { it.name == leaf.name }
            if (indexInConfigured >= 0) {
                configuredServices[indexInConfigured] = updatedLeaf
            } else {
                val indexInLinked = linkedServices.indexOfFirst { it.name == leaf.name }
                if (indexInLinked >= 0) {
                    linkedServices[indexInLinked] = updatedLeaf
                }
            }
        }
    }

    private fun readLeafRightsFromConfigFile(leafName: String, configFileName: String): Map<String, UserRight>? {
        return try {
            val configFile = if (configFileName == "docker.config") {
                findConfigFile("docker.config")
            } else {
                ConfigReader.loadConfig(configFileName)
            }
            
            if (configFile == null || !configFile.exists()) {
                return null
            }
            
            val lines = configFile.readLines()
            for (line in lines) {
                if (line.isBlank() || line.startsWith("#")) {
                    continue
                }
                
                val strSplt = line.split(';')
                if (strSplt.isNotEmpty() && strSplt[0].trim() == leafName) {
                    val allowedUsers = HashMap<String, UserRight>()
                    if (strSplt.size > 1) {
                        val userInfo = strSplt[1].split(',')
                        for (it in userInfo) {
                            val itemInfo = it.split(':')
                            if (itemInfo.size >= 2) {
                                val userName = itemInfo[0].trim()
                                val rights = UserRight.getFromString(itemInfo[1].trim())
                                allowedUsers[userName] = rights
                            }
                        }
                    }
                    return allowedUsers
                }
            }
            null
        } catch (e: Exception) {
            e.printStackTrace()
            null
        }
    }

    private fun updateLeafConfigFile(leafName: String, mergedRights: Map<String, UserRight>, configFileName: String) {
        try {
            val configFile = if (configFileName == "docker.config") {
                findConfigFile("docker.config")
            } else {
                ConfigReader.loadConfig(configFileName)
            }
            
            if (configFile == null || !configFile.exists()) {
                println("Config file $configFileName not found for leaf $leafName")
                return
            }
            
            // Читаем все строки файла
            val lines = configFile.readLines().toMutableList()
            
            // Ищем строку с нужным листом и обновляем её
            var found = false
            for (i in lines.indices) {
                val line = lines[i]
                // Пропускаем комментарии и пустые строки
                if (line.isBlank() || line.startsWith("#")) {
                    continue
                }
                
                val strSplt = line.split(';')
                if (strSplt.isNotEmpty() && strSplt[0].trim() == leafName) {
                    // Нашли строку с нужным листом - обновляем права
                    val rightsString = if (mergedRights.isNotEmpty()) {
                        mergedRights.map { (userName, rights) ->
                            val rightsStr = buildString {
                                if (rights.read) append('r')
                                if (rights.write) append('w')
                                if (rights.execute) append('x')
                            }
                            "$userName:$rightsStr"
                        }.joinToString(",")
                    } else {
                        ""
                    }
                    
                    // Формируем новую строку: leafName;user1:rwx,user2:rwx
                    lines[i] = if (rightsString.isNotEmpty()) {
                        "$leafName;$rightsString"
                    } else {
                        leafName // Если прав нет, сохраняем только имя листа
                    }
                    found = true
                    break
                }
            }
            
            if (!found) {
                println("Leaf $leafName not found in config file $configFileName")
                return
            }
            
            // Сохраняем обновлённый файл
            configFile.writeText(lines.joinToString("\n"))
        } catch (e: Exception) {
            e.printStackTrace()
            println("Error updating config file $configFileName for leaf $leafName: ${e.message}")
        }
    }
}