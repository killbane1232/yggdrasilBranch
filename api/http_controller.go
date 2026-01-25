package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"yggdrasil-branch/models"
	service "yggdrasil-branch/service"
	"yggdrasil-branch/ws"

	"github.com/gorilla/mux"
)

// HTTPController управляет HTTP API
type HTTPController struct {
	leafCollector *service.LeafCollector
	wsServer      *ws.WebSocketServer
}

// NewHTTPController создает новый HTTP контроллер
func NewHTTPController(leafCollector *service.LeafCollector, wsServer *ws.WebSocketServer) *HTTPController {
	return &HTTPController{
		leafCollector: leafCollector,
		wsServer:      wsServer,
	}
}

// RegisterRoutes регистрирует маршруты
func (hc *HTTPController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/leaf/connect", hc.HandleLeafConnect).Methods("POST")
	router.HandleFunc("/api/leaf/call/{leaf}", hc.HandleCallLeafMethod).Methods("POST")
	router.HandleFunc("/api/leaf/callback/{leaf}/{method}", hc.HandleCallback).Methods("POST")
	router.HandleFunc("/api/leaf/status", hc.HandleGetAllStatus).Methods("GET")
	router.HandleFunc("/api/leaf/status/{leaf}", hc.HandleGetStatus).Methods("GET")
	router.HandleFunc("/api/leaf/disconnect/{leaf}", hc.HandleDisconnect).Methods("DELETE")
	router.HandleFunc("/ws", hc.HandleWebSocket)
}

// HandleLeafConnect обрабатывает подключение листа
func (hc *HTTPController) HandleLeafConnect(w http.ResponseWriter, r *http.Request) {
	var leaf models.Leaf
	if err := json.NewDecoder(r.Body).Decode(&leaf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("HTTP: Connecting to /api/leaf/connect: %s\n", leaf.Name)
	hc.leafCollector.AddLinkedLeaf(&leaf)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// HandleCallLeafMethod обрабатывает вызов метода листа
func (hc *HTTPController) HandleCallLeafMethod(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leafName := vars["leaf"]

	var request models.LeafMethodMessage
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	leaves := hc.leafCollector.GetLeaves()
	var leaf *models.Leaf
	for _, l := range leaves {
		if l.Name == leafName {
			leaf = l
			break
		}
	}

	if leaf == nil {
		http.Error(w, "Leaf not found", http.StatusNotFound)
		return
	}

	// Если лист имеет URL, вызываем через HTTP
	if leaf.URL != "" {
		if controller, ok := leaf.Controller.(service.Controller); ok {
			result := controller.CallMethod(request.Method, request.Args)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(result))
			return
		}
	}

	// Иначе вызываем локально
	result := hc.leafCollector.CallServiceMethod(leafName, request.Method, request.Args)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

// HandleCallback обрабатывает callback от листа
func (hc *HTTPController) HandleCallback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leafName := vars["leaf"]
	method := vars["method"]

	var response string
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		// Пытаемся прочитать как обычный текст
		body := make([]byte, r.ContentLength)
		r.Body.Read(body)
		response = string(body)
	}

	// Обработка callback будет через RequestBuffer в будущем
	log.Printf("Callback from %s/%s: %s\n", leafName, method, response)
	w.Write([]byte("Callback processed"))
}

// HandleGetAllStatus возвращает статус всех листов
func (hc *HTTPController) HandleGetAllStatus(w http.ResponseWriter, r *http.Request) {
	leaves := hc.leafCollector.GetLeaves()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaves)
}

// HandleGetStatus возвращает статус конкретного листа
func (hc *HTTPController) HandleGetStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leafName := vars["leaf"]

	leaves := hc.leafCollector.GetLeaves()
	for _, leaf := range leaves {
		if leaf.Name == leafName {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(leaf)
			return
		}
	}

	http.Error(w, "Leaf not found", http.StatusNotFound)
}

// HandleDisconnect обрабатывает отключение листа
func (hc *HTTPController) HandleDisconnect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	leafName := vars["leaf"]

	// Удаление листа будет реализовано позже
	log.Printf("Disconnecting leaf: %s\n", leafName)
	w.Write([]byte(fmt.Sprintf("Leaf %s disconnected", leafName)))
}

// HandleWebSocket обрабатывает WebSocket подключения
func (hc *HTTPController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	hc.wsServer.HandleWebSocket(w, r, func(leaf *models.Leaf) *models.Leaf {
		log.Printf("Connecting to /ws: %s\n", leaf.Name)
		hc.leafCollector.AddLinkedLeaf(leaf)
		return leaf
	})
}
