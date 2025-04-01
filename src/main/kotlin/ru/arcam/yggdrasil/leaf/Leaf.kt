package ru.arcam.yggdrasil.leaf

import kotlinx.serialization.Serializable
import ru.arcam.yggdrasil.service.IController
import com.fasterxml.jackson.annotation.JsonIgnore
import kotlinx.serialization.Contextual

@Serializable
data class Leaf(
    val name: String,
    val status: String,
    val attachedBranch: String = "",
    val hooks: List<LeafHook>,
    @JsonIgnore
    @Contextual
    var controller: IController?
)
