package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	"yggdrasil-branch/models"
	"yggdrasil-branch/utils"
	"yggdrasil-branch/ws"
)

// LeafCollector управляет коллекцией листов
type LeafCollector struct {
	isWindows          bool
	leafStatus         map[string]string
	lock               sync.RWMutex
	linkedServices     []*models.Leaf
	configuredServices []*models.Leaf
	serviceName        string
	trunkConnection    *ws.TrunkConnection
}

// NewLeafCollector создает новый LeafCollector
func NewLeafCollector(trunkConnection *ws.TrunkConnection) *LeafCollector {
	isWindows := runtime.GOOS == "windows"
	return &LeafCollector{
		isWindows:          isWindows,
		leafStatus:         make(map[string]string),
		linkedServices:     make([]*models.Leaf, 0),
		configuredServices: make([]*models.Leaf, 0),
		serviceName:        utils.GetServiceName(),
		trunkConnection:    trunkConnection,
	}
}

// UpdateConfiguredServices обновляет список настроенных сервисов
func (lc *LeafCollector) UpdateConfiguredServices() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error in UpdateConfiguredServices: %v\n", r)
		}
	}()

	var confServicesBuffer []*models.Leaf

	lc.lock.Lock()
	// Сохраняем CustomController листы
	for _, leaf := range lc.configuredServices {
		if _, ok := leaf.Controller.(*CustomController); ok {
			confServicesBuffer = append(confServicesBuffer, leaf)
		}
	}
	lc.lock.Unlock()

	// Загружаем leaves.config
	configFile, err := utils.LoadConfig("leaves.config")
	if err == nil {
		scanner := bufio.NewScanner(configFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			leaf := lc.parseLeafFromConfigLine(line, lc.isWindows)
			if leaf != nil {
				if controller, ok := leaf.Controller.(Controller); ok {
					leaf.Status = controller.Status()
				}
				confServicesBuffer = append(confServicesBuffer, leaf)
			}
		}
		configFile.Close()
	}

	// Загружаем docker.config
	dockerConfigFile, err := utils.FindConfigFile("docker.config")
	if err == nil {
		scanner := bufio.NewScanner(dockerConfigFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			leaf := lc.parseLeafFromConfigLine(line, false) // Docker не зависит от ОС
			if leaf != nil {
				controller := NewDockerController(leaf)
				leaf.Controller = controller
				leaf.Status = controller.Status()
				confServicesBuffer = append(confServicesBuffer, leaf)
			}
		}
		dockerConfigFile.Close()
	}

	lc.lock.Lock()
	lc.configuredServices = confServicesBuffer
	lc.lock.Unlock()
}

// parseLeafFromConfigLine парсит строку конфигурации
func (lc *LeafCollector) parseLeafFromConfigLine(line string, isWindows bool) *models.Leaf {
	parts := strings.Split(line, ";")
	if len(parts) == 0 {
		return nil
	}

	leafName := strings.TrimSpace(parts[0])
	allowedUsers := make(map[string]models.UserRight)

	if len(parts) > 1 {
		userInfo := strings.Split(parts[1], ",")
		for _, userStr := range userInfo {
			userParts := strings.Split(userStr, ":")
			if len(userParts) >= 2 {
				userName := strings.TrimSpace(userParts[0])
				rights := models.GetUserRightFromString(strings.TrimSpace(userParts[1]))
				allowedUsers[userName] = rights
			}
		}
	}

	leaf := &models.Leaf{
		Name:           leafName,
		Status:         "UNAVAIVABLE",
		AttachedBranch: lc.serviceName,
		Hooks:          make([]models.LeafHook, 0),
		AllowedUsers:   allowedUsers,
	}

	var controller Controller
	if isWindows {
		controller = NewWindowsController(leaf)
	} else {
		controller = NewLinuxController(leaf)
	}
	leaf.Controller = controller

	return leaf
}

// ReportLeaves отправляет информацию о листах в Trunk
func (lc *LeafCollector) ReportLeaves() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error in ReportLeaves: %v\n", r)
		}
	}()

	var confServices []*models.Leaf
	lc.lock.RLock()
	confServices = append(confServices, lc.configuredServices...)
	confServices = append(confServices, lc.linkedServices...)
	lc.lock.RUnlock()

	// Загружаем глобальные права пользователей
	allowedUsers := lc.loadGlobalUserRights()

	// Преобразуем []*models.Leaf в []models.Leaf
	leaves := make([]models.Leaf, len(confServices))
	for i, leaf := range confServices {
		leaves[i] = *leaf
		// Убеждаемся, что AllowedUsers не nil
		if leaves[i].AllowedUsers == nil {
			leaves[i].AllowedUsers = make(map[string]models.UserRight)
		}
		// Убеждаемся, что Hooks не nil
		if leaves[i].Hooks == nil {
			leaves[i].Hooks = make([]models.LeafHook, 0)
		}
		// Обнуляем URL и Port при отправке BranchInfo (они не должны сериализоваться)
		leaves[i].URL = ""
		leaves[i].Port = 0
	}

	// Убеждаемся, что AllowedUsers не nil
	if allowedUsers == nil {
		allowedUsers = make(map[string]models.UserRight)
	}

	branchInfo := models.BranchInfo{
		ServiceName:  lc.serviceName,
		Leaves:       leaves,
		AllowedUsers: allowedUsers,
	}

	if lc.trunkConnection != nil {
		lc.trunkConnection.SendBranchInfo(branchInfo)
	}
}

// loadGlobalUserRights загружает глобальные права пользователей
func (lc *LeafCollector) loadGlobalUserRights() map[string]models.UserRight {
	allowedUsers := make(map[string]models.UserRight)

	userConfigFile, err := utils.LoadConfig("user.config")
	if err != nil {
		return allowedUsers
	}
	defer userConfigFile.Close()

	scanner := bufio.NewScanner(userConfigFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			userName := strings.TrimSpace(parts[0])
			rights := models.GetUserRightFromString(strings.TrimSpace(parts[1]))
			allowedUsers[userName] = rights
		}
	}

	return allowedUsers
}

// CallServiceMethod вызывает метод сервиса
func (lc *LeafCollector) CallServiceMethod(serviceName, method string, args []string) string {
	lc.lock.RLock()
	var leaf *models.Leaf
	for _, l := range lc.configuredServices {
		if l.Name == serviceName {
			leaf = l
			break
		}
	}
	if leaf == nil {
		for _, l := range lc.linkedServices {
			if l.Name == serviceName {
				leaf = l
				break
			}
		}
	}
	lc.lock.RUnlock()

	if leaf == nil {
		return "UNAVAIVABLE"
	}

	controller, ok := leaf.Controller.(Controller)
	if !ok {
		return "ERROR: Invalid controller"
	}

	switch method {
	case "STATUS":
		status := controller.Status()
		leaf.Status = status
		return status
	case "START":
		return controller.Start()
	case "STOP":
		return controller.Stop()
	case "RESTART":
		return controller.Restart()
	case "TAIL":
		return controller.Logs(nil)
	case "TAIL_N":
		return controller.Logs(args)
	case "RIGHTS":
		return lc.handleRightsCommand(serviceName, args)
	case "METHOD":
		// Для команды METHOD первый аргумент - это имя метода, остальные - аргументы
		if len(args) == 0 {
			return "ERROR: METHOD command requires at least method name"
		}
		fmt.Printf("METHOD command: serviceName=%s, args=%v\n", serviceName, "method", args)
		return controller.CallMethod("method", args)
	default:
		return controller.CallMethod(method, args)
	}
}

// handleRightsCommand обрабатывает команду RIGHTS
func (lc *LeafCollector) handleRightsCommand(leafName string, args []string) string {
	if len(args) == 0 {
		return "ERROR: No JSON provided"
	}

	jsonString := strings.Join(args, ":")
	fmt.Printf("Got message %s\n", jsonString)

	var rightsMap map[string]models.UserRight
	if err := json.Unmarshal([]byte(jsonString), &rightsMap); err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	lc.lock.Lock()
	defer lc.lock.Unlock()

	if leafName == "NULL" {
		lc.saveGlobalRights(rightsMap)
		return "OK: Global rights updated"
	}

	leaf := lc.findLeaf(leafName)
	if leaf == nil {
		return fmt.Sprintf("ERROR: Leaf %s not found", leafName)
	}

	lc.updateLeafRights(leaf, rightsMap)
	return fmt.Sprintf("OK: Rights updated for leaf %s", leafName)
}

// findLeaf находит лист по имени
func (lc *LeafCollector) findLeaf(name string) *models.Leaf {
	for _, leaf := range lc.configuredServices {
		if leaf.Name == name {
			return leaf
		}
	}
	for _, leaf := range lc.linkedServices {
		if leaf.Name == name {
			return leaf
		}
	}
	return nil
}

// saveGlobalRights сохраняет глобальные права
func (lc *LeafCollector) saveGlobalRights(newRightsMap map[string]models.UserRight) {
	configDir := "config"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Error creating config dir: %v\n", err)
		return
	}

	userConfigFile := filepath.Join(configDir, "user.config")

	existingRights := lc.loadGlobalUserRights()
	for k, v := range newRightsMap {
		existingRights[k] = v
	}

	var lines []string
	for userName, rights := range existingRights {
		lines = append(lines, fmt.Sprintf("%s:%s", userName, rights.ToString()))
	}

	if err := os.WriteFile(userConfigFile, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		fmt.Printf("Error writing user.config: %v\n", err)
	}
}

// updateLeafRights обновляет права листа
func (lc *LeafCollector) updateLeafRights(leaf *models.Leaf, newRightsMap map[string]models.UserRight) {
	existingRights := make(map[string]models.UserRight)
	for k, v := range leaf.AllowedUsers {
		existingRights[k] = v
	}

	for k, v := range newRightsMap {
		existingRights[k] = v
	}

	leaf.AllowedUsers = existingRights

	// Обновляем конфигурационный файл
	var configFileName string
	switch leaf.Controller.(type) {
	case *DockerController:
		configFileName = "docker.config"
	case *LinuxController:
		configFileName = "leaves.config"
	case *WindowsController:
		configFileName = "leaves.config"
	default:
		return // CustomController не сохраняем в файл
	}

	lc.updateLeafConfigFile(leaf.Name, existingRights, configFileName)
}

// updateLeafConfigFile обновляет конфигурационный файл листа
func (lc *LeafCollector) updateLeafConfigFile(leafName string, mergedRights map[string]models.UserRight, configFileName string) {
	configFile, err := utils.LoadConfig(configFileName)
	if err != nil {
		fmt.Printf("Config file %s not found for leaf %s\n", configFileName, leafName)
		return
	}
	defer configFile.Close()

	var lines []string
	scanner := bufio.NewScanner(configFile)
	found := false

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			lines = append(lines, line)
			continue
		}

		parts := strings.Split(trimmedLine, ";")
		if len(parts) > 0 && strings.TrimSpace(parts[0]) == leafName {
			rightsString := ""
			if len(mergedRights) > 0 {
				var rightsParts []string
				for userName, rights := range mergedRights {
					rightsParts = append(rightsParts, fmt.Sprintf("%s:%s", userName, rights.ToString()))
				}
				rightsString = strings.Join(rightsParts, ",")
			}

			if rightsString != "" {
				lines = append(lines, fmt.Sprintf("%s;%s", leafName, rightsString))
			} else {
				lines = append(lines, leafName)
			}
			found = true
		} else {
			lines = append(lines, line)
		}
	}

	if !found {
		fmt.Printf("Leaf %s not found in config file %s\n", leafName, configFileName)
		return
	}

	configPath := filepath.Join("config", configFileName)
	if err := os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		fmt.Printf("Error updating config file %s: %v\n", configFileName, err)
	}
}

// AddLinkedLeaf добавляет связанный лист
func (lc *LeafCollector) AddLinkedLeaf(leaf *models.Leaf) {
	lc.lock.Lock()
	defer lc.lock.Unlock()

	for _, l := range lc.linkedServices {
		if l.Name == leaf.Name {
			return
		}
	}

	// Убеждаемся, что AllowedUsers не nil
	if leaf.AllowedUsers == nil {
		leaf.AllowedUsers = make(map[string]models.UserRight)
	}
	// Убеждаемся, что Hooks не nil
	if leaf.Hooks == nil {
		leaf.Hooks = make([]models.LeafHook, 0)
	}

	leaf.AttachedBranch = lc.serviceName
	leaf.Controller = NewCustomController(leaf, lc)
	lc.linkedServices = append(lc.linkedServices, leaf)
}

// GetLeaves возвращает все листы
func (lc *LeafCollector) GetLeaves() []*models.Leaf {
	lc.lock.RLock()
	defer lc.lock.RUnlock()

	result := make([]*models.Leaf, 0)
	result = append(result, lc.configuredServices...)
	result = append(result, lc.linkedServices...)
	return result
}

// StartScheduler запускает планировщик задач
func (lc *LeafCollector) StartScheduler() {
	// Обновление настроенных сервисов каждые 5 секунд
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			lc.UpdateConfiguredServices()
		}
	}()

	// Отправка информации о листах каждые 5 секунд
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			lc.ReportLeaves()
		}
	}()
}
