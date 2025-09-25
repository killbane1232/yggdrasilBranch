package ru.arcam.yggdrasil.leaf

import kotlinx.serialization.Serializable
import ru.arcam.yggdrasil.service.IController
import com.fasterxml.jackson.annotation.JsonIgnore
import com.fasterxml.jackson.annotation.JsonProperty
import kotlinx.serialization.Contextual
import lombok.AllArgsConstructor
import lombok.Builder
import lombok.NoArgsConstructor
import ru.arcam.yggdrasil.user.UserRight

@Serializable
@Builder
@NoArgsConstructor
@AllArgsConstructor
data class Leaf(
    @JsonProperty
    val name: String,
    @JsonProperty
    var status: String,
    @JsonProperty
    var attachedBranch: String = "",
    @JsonProperty
    val hooks: List<LeafHook>,
    @JsonProperty
    val allowedUsers : Map<String, UserRight> = hashMapOf(),
    @JsonIgnore
    @Contextual
    var controller: IController? = null,
    @JsonIgnore
    var port: Int = 0,
    @JsonIgnore
    var url: String = ""
)
