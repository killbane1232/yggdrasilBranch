package ru.arcam.yggdrasil.user

import com.fasterxml.jackson.annotation.JsonProperty
import kotlinx.serialization.Serializable
import lombok.Builder

@Serializable
@Builder
data class UserRight(
    @JsonProperty
    val read: Boolean,
    @JsonProperty
    val write: Boolean,
    @JsonProperty
    val execute: Boolean
) {
    companion object {
        fun getFromString(right: String) : UserRight {
            return UserRight(
                right.contains('r', true),
                right.contains('w', true),
                right.contains('x', true)
            )
        }
    }
}