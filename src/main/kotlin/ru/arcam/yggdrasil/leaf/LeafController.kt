package ru.arcam.yggdrasil.leaf

import jakarta.annotation.PostConstruct
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.messaging.handler.annotation.DestinationVariable
import org.springframework.messaging.handler.annotation.MessageMapping
import org.springframework.messaging.handler.annotation.Payload
import org.springframework.messaging.handler.annotation.SendTo
import org.springframework.messaging.simp.SimpMessagingTemplate
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Controller
import ru.arcam.yggdrasil.utils.NameResolver
import ru.arcam.yggdrasil.ws.RequestBuffer


@Controller
class LeafController(private val messagingTemplate: SimpMessagingTemplate) {
    @Autowired
    private lateinit var leafCollector: LeafCollector

    @Autowired
    private lateinit var requestBuffer: RequestBuffer

    @PostConstruct
    fun afterCreate() {
        leafController = this
    }

    @MessageMapping("/connect")
    @SendTo("/topic/leaf-status")
    fun handleLeafConnection(@Payload leaf: Leaf): Leaf {
        synchronized(leafCollector.lock) {
            leaf.attachedBranch = NameResolver.name
            leafCollector.configuredServices.add(leaf)
        }
        return leaf
    }

    @MessageMapping("/callback/{leaf}/{method}")
    fun handleCallback(
        @DestinationVariable leaf: String,
        @DestinationVariable method: String,
        response: String
    ) {
        requestBuffer.handleResponse(leaf, method, response)
    }

    @Scheduled(fixedRate = 60000) // Каждую минуту
    fun cleanupRequests() {
        requestBuffer.cleanup()
    }

    @MessageMapping("/leaf/{leaf}")
    @SendTo("/topic/leaf/{leaf}")
    fun callLeafMethod(@DestinationVariable leaf: String, method: String, args: List<String>): String {
        println("Sending message to /topic/leaf/$leaf: $method")

        messagingTemplate.convertAndSend("/topic/leaf/$leaf", LeafMethodMessage(method, args))
        return requestBuffer.addRequest(leaf, method).get()
    }

    @MessageMapping("/topic/message/{leaf}")
    fun callBranchMethod(@DestinationVariable leaf: String, method: String, args: List<String>): String {
        println("Sending message to /topic/leaf/$leaf: $method")

        messagingTemplate.convertAndSend("/topic/leaf/$leaf", LeafMethodMessage(method, args))
        return requestBuffer.addRequest(leaf, method).get()
    }

    private data class LeafMethodMessage(
        val method: String,
        val args: List<String>
    )

    companion object {
        lateinit var leafController: LeafController
    }
}
