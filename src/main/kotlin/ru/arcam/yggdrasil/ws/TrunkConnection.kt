package ru.arcam.yggdrasil.ws

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.messaging.converter.MappingJackson2MessageConverter
import org.springframework.messaging.simp.stomp.*
import org.springframework.scheduling.annotation.Scheduled
import org.springframework.stereotype.Service
import org.springframework.web.socket.WebSocketHttpHeaders
import org.springframework.web.socket.client.standard.StandardWebSocketClient
import org.springframework.web.socket.messaging.WebSocketStompClient
import ru.arcam.yggdrasil.leaf.Leaf
import ru.arcam.yggdrasil.leaf.LeafCollector
import java.lang.reflect.Type
import kotlin.reflect.javaType
import kotlin.reflect.typeOf

@Service
class TrunkConnection {
    @Autowired
    lateinit var leafCollector: LeafCollector
    
    private val config = TrunkConfig()

    @Scheduled(fixedRate = 1000)
    fun restoreConnection() {
        if (WSClient == null || !WSClient!!.isConnected) {
            val stompClient = WebSocketStompClient(StandardWebSocketClient())
            stompClient.messageConverter = MappingJackson2MessageConverter()

            val sessionHandler: StompSessionHandler = MyStompSessionHandler(leafCollector)
            stompClient.start()
            config.loadConfig()
            println("Connecting to ${config.getWebSocketUrl()}")
            val websocketheaders = WebSocketHttpHeaders()
            val headers = StompHeaders()

            val res = stompClient.connectAsync(config.getWebSocketUrl(), websocketheaders, headers, sessionHandler)
            res.whenComplete { x, y ->
                if (x != null) {
                    WSClient = x
                }
                else
                    println(y)
                println("Connected: " + (x != null))
            }
        }
    }

    companion object {
        var WSClient: StompSession? = null
    }

    class MyStompSessionHandler(val leafCollector: LeafCollector) : StompSessionHandlerAdapter() {
        val serviceName: String = java.net.InetAddress.getLocalHost().hostName
        @OptIn(ExperimentalStdlibApi::class)
        override fun getPayloadType(headers: StompHeaders): Type {
            return typeOf<String>().javaType
        }

        override fun handleFrame(headers: StompHeaders, payload: Any?) {
            println("Received message: $payload from ${headers.destination}")
            val payloadSplit = payload.toString().split(':')
            val args = listOf<String>()
            if (payloadSplit.size > 2)
                payloadSplit.subList(2, payloadSplit.size - 1)
            val result = leafCollector.callServiceMethod(payloadSplit[1], payloadSplit[0], args)
            WSClient!!.send("/app/callback/$serviceName", result);
        }

        override fun afterConnected(session: StompSession, connectedHeaders: StompHeaders) {
            println("Opening websocket...");
            println("Session is connected: " + session.isConnected() + "\n")
            session.subscribe("/topic/message/$serviceName", this)
        }

        override fun handleException(
            session: StompSession,
            command: StompCommand?,
            headers: StompHeaders,
            payload: ByteArray,
            exception: Throwable
        ) {
            throw exception
        }

        override fun handleTransportError(session: StompSession, exception: Throwable) {
            throw exception
        }
    }
}