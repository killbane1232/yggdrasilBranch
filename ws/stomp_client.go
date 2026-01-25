package ws

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

// STOMPFrame представляет STOMP фрейм
type STOMPFrame struct {
	Command string
	Headers map[string]string
	Body    string
}

// Encode кодирует STOMP фрейм в строку
func (f *STOMPFrame) Encode() string {
	var builder strings.Builder

	// Команда
	builder.WriteString(f.Command)
	builder.WriteString("\n")

	// Заголовки
	for key, value := range f.Headers {
		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(value)
		builder.WriteString("\n")
	}

	// Пустая строка перед телом
	builder.WriteString("\n")

	// Тело
	if f.Body != "" {
		builder.WriteString(f.Body)
	}

	// NULL символ в конце
	builder.WriteString("\x00")

	return builder.String()
}

// DecodeSTOMPFrame декодирует STOMP фрейм из WebSocket сообщения
func DecodeSTOMPFrame(message []byte) (*STOMPFrame, error) {
	frame := &STOMPFrame{
		Headers: make(map[string]string),
	}

	// Разделяем на строки
	parts := strings.Split(string(message), "\n")
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty STOMP frame")
	}

	// Первая строка - команда
	frame.Command = strings.TrimSpace(parts[0])

	// Читаем заголовки
	headerEnd := 1
	for i := 1; i < len(parts); i++ {
		line := strings.TrimSpace(parts[i])
		if line == "" {
			headerEnd = i + 1
			break // Пустая строка означает конец заголовков
		}

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			frame.Headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}

	// Остальное - тело (до NULL символа)
	if headerEnd < len(parts) {
		bodyParts := parts[headerEnd:]
		bodyStr := strings.Join(bodyParts, "\n")
		// Удаляем NULL символ в конце
		bodyStr = strings.TrimSuffix(bodyStr, "\x00")
		frame.Body = bodyStr
	}

	return frame, nil
}

// SendSTOMP отправляет STOMP фрейм через WebSocket
func SendSTOMP(conn *websocket.Conn, command string, headers map[string]string, body interface{}) error {
	var bodyStr string
	if body != nil {
		if str, ok := body.(string); ok {
			bodyStr = str
		} else {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal body: %w", err)
			}
			bodyStr = string(jsonData)
		}
	}

	frame := &STOMPFrame{
		Command: command,
		Headers: headers,
		Body:    bodyStr,
	}

	return conn.WriteMessage(websocket.TextMessage, []byte(frame.Encode()))
}
