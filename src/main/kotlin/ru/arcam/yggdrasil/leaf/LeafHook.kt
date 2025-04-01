package ru.arcam.yggdrasil.leaf

import kotlinx.serialization.Serializable

@Serializable
data class LeafHook (
    val name: String,
    val hookFields: Map<String, HookType>
)

enum class HookType (className: String) {
    STRING(String::class.java.name)
}