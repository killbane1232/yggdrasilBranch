package ru.arcam.yggdrasil.branch

import kotlinx.serialization.Serializable
import ru.arcam.yggdrasil.leaf.Leaf
import ru.arcam.yggdrasil.user.UserRight

@Serializable
data class BranchInfo(
    val serviceName : String,
    val leaves : List<Leaf>,
    val allowedUsers : Map<String, UserRight> = hashMapOf()
)
