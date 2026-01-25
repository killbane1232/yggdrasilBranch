package ws

import (
	"fmt"
	"log"
	"net/http"
	"yggdrasil-branch/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все источники
	},
}

// WebSocketServer управляет WebSocket сервером для подключения листов
type WebSocketServer struct {
	clients map[string]*websocket.Conn
}

// NewWebSocketServer создает новый WebSocket сервер
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		clients: make(map[string]*websocket.Conn),
	}
}

// HandleWebSocket обрабатывает WebSocket подключения
func (wss *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request, onLeafConnect func(*models.Leaf) *models.Leaf) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}
	defer conn.Close()

	for {
		var leaf models.Leaf
		if err := conn.ReadJSON(&leaf); err != nil {
			log.Printf("Error reading message: %v\n", err)
			break
		}

		// Обрабатываем подключение листа
		if leaf.Name != "" {
			updatedLeaf := onLeafConnect(&leaf)
			if err := conn.WriteJSON(updatedLeaf); err != nil {
				log.Printf("Error writing message: %v\n", err)
				break
			}
			wss.clients[leaf.Name] = conn
		}
	}
}

// SendToLeaf отправляет сообщение листу
func (wss *WebSocketServer) SendToLeaf(leafName string, message models.LeafMethodMessage) error {
	conn, ok := wss.clients[leafName]
	if !ok {
		return fmt.Errorf("leaf %s not connected", leafName)
	}

	return conn.WriteJSON(message)
}
