package ru.arcam.yggdrasil.branch

data class BranchMessage(
    val branchName: String,
    val leafName: String,
    val messageData: Map<String, Any>
)
