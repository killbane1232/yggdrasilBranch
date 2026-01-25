package service

import "yggdrasil-branch/models"

// Controller интерфейс для управления сервисами
type Controller interface {
	Stop() string
	Start() string
	Restart() string
	Status() string
	Logs(args []string) string
	CallMethod(method string, args []string) string
	GetLeaf() *models.Leaf
}
