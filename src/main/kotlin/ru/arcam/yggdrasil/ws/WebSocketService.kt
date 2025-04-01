package ru.arcam.yggdrasil.ws

import jakarta.annotation.PostConstruct
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.messaging.handler.annotation.DestinationVariable
import org.springframework.messaging.handler.annotation.MessageMapping
import org.springframework.messaging.handler.annotation.Payload
import org.springframework.messaging.handler.annotation.SendTo
import org.springframework.messaging.simp.SimpMessagingTemplate
import org.springframework.stereotype.Controller
import org.springframework.stereotype.Service


@Controller
class WebSocketService(private val messagingTemplate: SimpMessagingTemplate) {
    @PostConstruct
    fun afterCreate() {
        wsService = this
    }

    @SendTo("/topic/callback/{username}")
    fun processMessage(@DestinationVariable username: String, @Payload data: String) {
        println("Sending message to /topic/message/$username: $data")
        messagingTemplate.convertAndSend("/topic/callback/$username", data)
    }

    companion object {
        lateinit var wsService: WebSocketService
    }
}