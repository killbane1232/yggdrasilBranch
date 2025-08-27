package ru.arcam.yggdrasil.leaf

data class LeafMethodMessage(
    val method: String,
    val args: List<String>
)