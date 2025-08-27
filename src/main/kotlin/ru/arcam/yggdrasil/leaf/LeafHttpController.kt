package ru.arcam.yggdrasil.leaf

import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.annotation.PostConstruct
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.ResponseEntity
import org.springframework.messaging.simp.SimpMessagingTemplate
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.web.bind.annotation.*
import ru.arcam.yggdrasil.service.CustomController
import ru.arcam.yggdrasil.utils.NameResolver
import ru.arcam.yggdrasil.ws.RequestBuffer
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

    @Autowired
    private lateinit var requestBuffer: RequestBuffer

    @PostConstruct
    fun afterCreate() {
        leafHttpController = this
    }

    /**
     * Подключение листа к ветке
     */
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

    /**
     * Обработка callback от листа
     */
    @PostMapping("/callback/{leaf}")
    fun handleCallback(
        @PathVariable leaf: String,
        @PathVariable method: String,
        @RequestBody response: String
    ): ResponseEntity<String> {
        requestBuffer.handleResponse(leaf, method, response)
        return ResponseEntity.ok("Callback processed")
    }

    fun callLeafMethod(
        leaf: Leaf,
        request: LeafMethodMessage
    ): ResponseEntity<String> {
        println("HTTP: Calling method on leaf $leaf: ${request.method}")

        val objectMapper = ObjectMapper()
        val requestBody: String = objectMapper
            .writeValueAsString(request)
        // Отправляем сообщение через WebSocket для совместимости
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

    /**
     * Получение статуса всех подключенных листов
     */
    @GetMapping("/status")
    fun getLeafStatus(): ResponseEntity<List<Leaf>> {
        return ResponseEntity.ok(leafCollector.linkedServices)
    }

    /**
     * Получение статуса конкретного листа
     */
    @GetMapping("/status/{leaf}")
    fun getLeafStatus(@PathVariable leaf: String): ResponseEntity<Leaf?> {
        val foundLeaf = leafCollector.linkedServices.firstOrNull { it.name == leaf }
        return if (foundLeaf != null) {
            ResponseEntity.ok(foundLeaf)
        } else {
            ResponseEntity.notFound().build()
        }
    }

    /**
     * Отключение листа
     */
    @DeleteMapping("/disconnect/{leaf}")
    fun disconnectLeaf(@PathVariable leaf: String): ResponseEntity<String> {
        synchronized(leafCollector.lock) {
            leafCollector.linkedServices.removeAll { it.name == leaf }
        }
        return ResponseEntity.ok("Leaf $leaf disconnected")
    }

    companion object {
        lateinit var leafHttpController: LeafHttpController
    }
}
