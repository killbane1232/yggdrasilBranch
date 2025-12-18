package ru.arcam.yggdrasil.user

import com.fasterxml.jackson.annotation.JsonProperty
import kotlinx.serialization.Serializable
import lombok.Builder
import lombok.NoArgsConstructor

@Serializable
@Builder
data class UserRight(
    @JsonProperty
    var read: Boolean,
    @JsonProperty
    var write: Boolean,
    @JsonProperty
    var execute: Boolean
) {
    constructor() : this(false, false, false)
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