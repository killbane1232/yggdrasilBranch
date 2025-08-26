package ru.arcam.yggdrasil.leaf

import com.fasterxml.jackson.annotation.JsonProperty
import kotlinx.serialization.Serializable

@Serializable
data class LeafHook (
    @JsonProperty
    val name: String,
    @JsonProperty
    val hookFields: Map<String, HookType>
)

enum class HookType (className: String) {
    STRING(String::class.java.name)
}