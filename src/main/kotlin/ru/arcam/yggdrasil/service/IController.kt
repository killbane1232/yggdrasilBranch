package ru.arcam.yggdrasil.service
import ru.arcam.yggdrasil.leaf.Leaf

abstract class IController(val leaf: Leaf) {
    abstract fun stop(): String
    abstract fun start(): String
    abstract fun restart(): String
    abstract fun status(): String
    abstract fun logs(args: List<String>? = null): String
    abstract fun callMethod(method: String, args: List<String>): String
}
