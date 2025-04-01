package ru.arcam.yggdrasil.branch

import kotlinx.serialization.Serializable
import ru.arcam.yggdrasil.leaf.Leaf

@Serializable
data class BranchInfo(
    val serviceName : String,
    val leaves : List<Leaf>
)
