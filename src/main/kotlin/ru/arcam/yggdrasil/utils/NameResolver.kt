package ru.arcam.yggdrasil.utils

class NameResolver {
    companion object {
        val name: String = System.getenv("BRANCHNAME")?.takeIf { it.isNotBlank() } 
            ?: java.net.InetAddress.getLocalHost().hostName
    }
}