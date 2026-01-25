package ws

import (
	"fmt"
	"log"
	"strings"
	"time"
	"yggdrasil-branch/models"
	"yggdrasil-branch/utils"

	"github.com/gorilla/websocket"
)

// TrunkConnection управляет подключением к Trunk через WebSocket
type TrunkConnection struct {
	conn           *websocket.Conn
	config         *utils.TrunkConfig
	serviceName    string
	messageHandler func(method, leafName string, args []string) string
	onConnected    func() // Callback при успешном подключении
}

// NewTrunkConnection создает новое подключение к Trunk
func NewTrunkConnection(serviceName string) *TrunkConnection {
	return &TrunkConnection{
		config:      utils.NewTrunkConfig(),
		serviceName: serviceName,
	}
}

// Connect подключается к Trunk
func (tc *TrunkConnection) Connect() error {
	url := tc.config.GetWebSocketURL()
	dialer := websocket.Dialer{
		HandshakeTimeout: time.Duration(tc.config.GetTimeout()) * time.Millisecond,
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	tc.conn = conn
	fmt.Printf("Connected to %s\n", url)

	// Отправляем STOMP CONNECT
	connectHeaders := map[string]string{
		"accept-version": "1.1,1.0",
		"heart-beat":     "10000,10000",
	}
	if err := SendSTOMP(conn, "CONNECT", connectHeaders, nil); err != nil {
		conn.Close()
		return fmt.Errorf("failed to send CONNECT: %w", err)
	}

	// Ждем CONNECTED фрейм
	_, message, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to read CONNECTED: %w", err)
	}

	connectedFrame, err := DecodeSTOMPFrame(message)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to decode CONNECTED: %w", err)
	}

	if connectedFrame.Command != "CONNECTED" {
		conn.Close()
		return fmt.Errorf("expected CONNECTED, got %s", connectedFrame.Command)
	}

	fmt.Printf("STOMP connected, session: %s\n", connectedFrame.Headers["session"])

	// Подписываемся на топик
	subscribeHeaders := map[string]string{
		"id":          "sub-0",
		"destination": fmt.Sprintf("/topic/message/%s", tc.serviceName),
	}
	if err := SendSTOMP(conn, "SUBSCRIBE", subscribeHeaders, nil); err != nil {
		conn.Close()
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	fmt.Printf("Subscribed to /topic/message/%s\n", tc.serviceName)

	// Всегда запускаем обработчик сообщений
	go tc.handleMessages()

	// Вызываем callback при успешном подключении
	if tc.onConnected != nil {
		tc.onConnected()
	}

	return nil
}

// IsConnected проверяет, подключен ли клиент
func (tc *TrunkConnection) IsConnected() bool {
	return tc.conn != nil
}

// SendBranchInfo отправляет информацию о ветке
func (tc *TrunkConnection) SendBranchInfo(branchInfo models.BranchInfo) {
	if !tc.IsConnected() {
		log.Printf("Cannot send branch info: not connected\n")
		return
	}

	headers := map[string]string{
		"destination":  "/app/assing",
		"content-type": "application/json",
	}

	log.Printf("Sending branch info: serviceName=%s, leavesCount=%d\n", branchInfo.ServiceName, len(branchInfo.Leaves))
	if err := SendSTOMP(tc.conn, "SEND", headers, branchInfo); err != nil {
		log.Printf("Error sending branch info: %v\n", err)
		tc.conn = nil
	} else {
		log.Printf("Branch info sent successfully\n")
	}
}

// SendCallback отправляет callback
func (tc *TrunkConnection) SendCallback(result string) {
	if !tc.IsConnected() {
		return
	}

	headers := map[string]string{
		"destination":  fmt.Sprintf("/app/callback/%s", tc.serviceName),
		"content-type": "text/plain",
	}

	if err := SendSTOMP(tc.conn, "SEND", headers, result); err != nil {
		log.Printf("Error sending callback: %v\n", err)
		tc.conn = nil
	}
}

// handleMessages обрабатывает входящие сообщения
func (tc *TrunkConnection) handleMessages() {
	log.Printf("Message handler started\n")
	for {
		if tc.conn == nil {
			log.Printf("Message handler: connection is nil, exiting\n")
			return
		}

		_, message, err := tc.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading WebSocket message: %v\n", err)
			tc.conn = nil
			return
		}

		frame, err := DecodeSTOMPFrame(message)
		if err != nil {
			log.Printf("Error decoding STOMP frame: %v, raw message: %s\n", err, string(message))
			continue
		}

		log.Printf("Received STOMP frame: command=%s\n", frame.Command)

		if frame.Command == "MESSAGE" {
			log.Printf("STOMP MESSAGE received: destination=%s, body=%s\n", frame.Headers["destination"], frame.Body)
			// Парсим сообщение в формате "method:leafName:arg1:arg2:..."
			parts := strings.Split(strings.Trim(frame.Body, "\""), ":")
			if len(parts) < 2 {
				log.Printf("Invalid message format, expected 'method:leafName:args', got: %s\n", frame.Body)
				continue
			}

			method := parts[0]
			leafName := parts[1]
			args := parts[2:]

			log.Printf("Processing message: method=%s, leaf=%s, args=%v\n", method, leafName, args)

			if tc.messageHandler != nil {
				result := tc.messageHandler(method, leafName, args)
				log.Printf("Handler result: %s\n", result)
				tc.SendCallback(result)
			} else {
				log.Printf("No message handler set, ignoring message\n")
			}
		} else if frame.Command == "ERROR" {
			log.Printf("STOMP ERROR: %s\n", frame.Body)
		} else if frame.Command == "CONNECTED" {
			log.Printf("STOMP CONNECTED received (duplicate)\n")
		} else {
			log.Printf("Unknown STOMP command: %s\n", frame.Command)
		}
	}
}

// SetOnConnected устанавливает callback при успешном подключении
func (tc *TrunkConnection) SetOnConnected(callback func()) {
	tc.onConnected = callback
}

// RestoreConnection восстанавливает подключение
func (tc *TrunkConnection) RestoreConnection(handler func(method, leafName string, args []string) string) {
	// Устанавливаем handler перед подключением
	tc.messageHandler = handler

	if tc.IsConnected() {
		return
	}

	if err := tc.Connect(); err != nil {
		log.Printf("Failed to restore connection: %v\n", err)
		return
	}

	// Запускаем цикл переподключения
	go func() {
		for {
			if !tc.IsConnected() {
				time.Sleep(1 * time.Second)
				if err := tc.Connect(); err != nil {
					log.Printf("Failed to reconnect: %v\n", err)
				} else {
					// После переподключения handler уже установлен
					tc.messageHandler = handler
				}
			} else {
				time.Sleep(5 * time.Second) // Проверяем подключение каждые 5 секунд
			}
		}
	}()
}
