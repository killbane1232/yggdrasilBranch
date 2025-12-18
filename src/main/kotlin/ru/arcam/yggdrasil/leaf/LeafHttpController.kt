package ru.arcam.yggdrasil.leaf

import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.annotation.PostConstruct
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.ResponseEntity
import org.springframework.messaging.simp.SimpMessagingTemplate
import org.springframework.web.bind.annotation.*
import ru.arcam.yggdrasil.service.CustomController
import ru.arcam.yggdrasil.utils.NameResolver
import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpRequest.BodyPublishers
import java.net.http.HttpResponse

@RestController
@RequestMapping("/api/leaf")
class LeafHttpController(private val messagingTemplate: SimpMessagingTemplate) {
    @Autowired
    private lateinit var leafCollector: LeafCollector

    @PostConstruct
    fun afterCreate() {
        leafHttpController = this
    }

    @PostMapping("/connect")
    fun handleLeafConnection(@RequestBody leaf: Leaf): ResponseEntity<String> {
        println("HTTP: Connecting to /api/leaf/connect: ${leaf.name}")
        synchronized(leafCollector.lock) {
            leaf.attachedBranch = NameResolver.name
            leaf.controller = CustomController(leaf, leafCollector)
            if (leafCollector.linkedServices.firstOrNull { x -> x.name == leaf.name } == null)
                leafCollector.linkedServices.add(leaf)
        }
        return ResponseEntity.ok("ok")
    }

    fun callLeafMethod(
        leaf: Leaf,
        request: LeafMethodMessage
    ): ResponseEntity<String> {
        println("HTTP: Calling method on leaf $leaf: ${request.method}")

        val objectMapper = ObjectMapper()
        val requestBody: String = objectMapper
            .writeValueAsString(request)

        val url = leaf.url + "/api/leaf/invoke"
        val client = HttpClient.newBuilder().build();
        val requestt = HttpRequest.newBuilder()
            .uri(URI.create(url))
            .POST(BodyPublishers.ofString(requestBody))
            .setHeader("Content-Type", "application/json")
            .build();

        val response = client.send(requestt, HttpResponse.BodyHandlers.ofString())
        return ResponseEntity.ok(response.body())
    }

    companion object {
        lateinit var leafHttpController: LeafHttpController
    }
}
