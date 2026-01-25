package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"yggdrasil-branch/api"
	"yggdrasil-branch/service"
	"yggdrasil-branch/utils"
	"yggdrasil-branch/ws"

	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("port", "8081", "Порт для HTTP сервера")
	flag.Parse()

	serviceName := utils.GetServiceName()
	fmt.Printf("Starting Yggdrasil Branch service: %s\n", serviceName)

	// Создаем TrunkConnection
	trunkConnection := ws.NewTrunkConnection(serviceName)

	// Создаем LeafCollector
	leafCollector := service.NewLeafCollector(trunkConnection)

	// Устанавливаем callback при подключении - сразу отправляем BranchInfo
	trunkConnection.SetOnConnected(func() {
		fmt.Printf("Connection established, sending initial branch info...\n")
		// Даем немного времени на завершение подписки
		time.Sleep(500 * time.Millisecond)
		leafCollector.ReportLeaves()
	})

	// Запускаем планировщик задач
	leafCollector.StartScheduler()

	// Запускаем восстановление подключения к Trunk
	go func() {
		for {
			if !trunkConnection.IsConnected() {
				trunkConnection.RestoreConnection(func(method, leafName string, args []string) string {
					return leafCollector.CallServiceMethod(leafName, method, args)
				})
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Создаем WebSocket сервер
	wsServer := ws.NewWebSocketServer()

	// Создаем HTTP контроллер
	httpController := api.NewHTTPController(leafCollector, wsServer)

	// Настраиваем маршруты
	router := mux.NewRouter()
	httpController.RegisterRoutes(router)

	// Запускаем HTTP сервер
	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("HTTP server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
