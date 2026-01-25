package utils

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// TrunkConfig представляет конфигурацию подключения к Trunk
type TrunkConfig struct {
	Host    string
	Port    int
	Path    string
	Timeout int64
}

// NewTrunkConfig создает новую конфигурацию и загружает её из файла
func NewTrunkConfig() *TrunkConfig {
	config := &TrunkConfig{
		Host:    "localhost",
		Port:    8080,
		Path:    "/ws",
		Timeout: 5000,
	}
	config.LoadConfig()
	return config
}

// LoadConfig загружает конфигурацию из файла
func (tc *TrunkConfig) LoadConfig() {
	configFile, err := LoadConfig("websocket.config")
	if err != nil {
		fmt.Printf("Warning: could not load websocket.config: %v\n", err)
		return
	}
	defer configFile.Close()
	
	scanner := bufio.NewScanner(configFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		switch key {
		case "websocket.host":
			tc.Host = value
		case "websocket.port":
			if port, err := strconv.Atoi(value); err == nil {
				tc.Port = port
			}
		case "websocket.path":
			tc.Path = value
		case "websocket.timeout":
			if timeout, err := strconv.ParseInt(value, 10, 64); err == nil {
				tc.Timeout = timeout
			}
		}
	}
	
	fmt.Printf("Config after reading: %s\n", tc.GetWebSocketURL())
}

// GetWebSocketURL возвращает URL WebSocket подключения
func (tc *TrunkConfig) GetWebSocketURL() string {
	return fmt.Sprintf("ws://%s:%d%s", tc.Host, tc.Port, tc.Path)
}

// GetTimeout возвращает таймаут подключения
func (tc *TrunkConfig) GetTimeout() int64 {
	return tc.Timeout
}
