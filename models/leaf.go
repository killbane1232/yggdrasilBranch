package models

// Leaf представляет лист (сервис)
// Controller должен быть интерфейсом service.Controller, но мы не можем импортировать service в models
// Поэтому используем interface{} и приведение типов в service пакете
type Leaf struct {
	Name           string               `json:"name"`
	Status         string               `json:"status"`
	AttachedBranch string               `json:"attachedBranch"`
	Hooks          []LeafHook           `json:"hooks"`
	AllowedUsers   map[string]UserRight `json:"allowedUsers"`
	Controller     interface{}          `json:"-"`
	Port           int                  `json:"port,omitempty"` // Читается при /connect, но не отправляется в BranchInfo
	URL            string               `json:"url,omitempty"`  // Читается при /connect, но не отправляется в BranchInfo
}

// LeafHook представляет хук листа
type LeafHook struct {
	Name       string            `json:"name"`
	HookFields map[string]string `json:"hookFields"`
}

// LeafMethodMessage представляет сообщение для вызова метода листа
type LeafMethodMessage struct {
	Method string   `json:"method"`
	Args   []string `json:"args"`
}
