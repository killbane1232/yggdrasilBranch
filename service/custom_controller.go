package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"yggdrasil-branch/models"
)

// CustomController управляет кастомными листами через HTTP
type CustomController struct {
	leaf          *models.Leaf
	leafCollector *LeafCollector
}

// NewCustomController создает новый CustomController
func NewCustomController(leaf *models.Leaf, collector *LeafCollector) *CustomController {
	return &CustomController{
		leaf:          leaf,
		leafCollector: collector,
	}
}

// GetLeaf возвращает лист
func (cc *CustomController) GetLeaf() *models.Leaf {
	return cc.leaf
}

// Stop останавливает сервис
func (cc *CustomController) Stop() string {
	cc.CallMethod("stop", []string{})
	return "OK"
}

// Start запускает сервис
func (cc *CustomController) Start() string {
	cc.CallMethod("start", []string{})
	return "OK"
}

// Restart перезапускает сервис
func (cc *CustomController) Restart() string {
	cc.CallMethod("restart", []string{})
	return "OK"
}

// Status возвращает статус сервиса
func (cc *CustomController) Status() string {
	status := cc.CallMethod("status", []string{})
	if status == "ERROR" {
		return "UNAVAIVABLE"
	}
	return status
}

// Logs возвращает логи сервиса
func (cc *CustomController) Logs(args []string) string {
	return "UNAVAIVABLE"
}

// CallMethod вызывает метод через HTTP
func (cc *CustomController) CallMethod(method string, args []string) string {
	if cc.leaf.URL == "" {
		log.Printf("HTTP: Empty URL: %s\n", cc.leaf.Name)
		return "ERROR"
	}

	request := models.LeafMethodMessage{
		Method: method,
		Args:   args,
	}

	log.Printf("HTTP: Calling method on leaf %s: method=%s, args=%v\n", cc.leaf.Name, method, args)

	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("HTTP: Error creating request: %s\n", err)
		return "ERROR"
	}

	log.Printf("HTTP: Request JSON: %s\n", string(jsonData))

	url := cc.leaf.URL + "/api/leaf/invoke"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("HTTP: Error creating request: %s\n", err)
		return "ERROR"
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP: Error sending request: %s\n", err)
		return "ERROR"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("HTTP: Error reading response: %s\n", err)
		return "ERROR"
	}

	return string(body)
}
