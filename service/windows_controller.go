package service

import (
	"bufio"
	"os/exec"
	"regexp"
	"strings"
	"yggdrasil-branch/models"
)

// WindowsController управляет Windows сервисами через sc
type WindowsController struct {
	leaf *models.Leaf
}

// NewWindowsController создает новый WindowsController
func NewWindowsController(leaf *models.Leaf) *WindowsController {
	return &WindowsController{leaf: leaf}
}

// GetLeaf возвращает лист
func (wc *WindowsController) GetLeaf() *models.Leaf {
	return wc.leaf
}

// Stop останавливает сервис
func (wc *WindowsController) Stop() string {
	cmd := exec.Command("sc", "stop", wc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Start запускает сервис
func (wc *WindowsController) Start() string {
	cmd := exec.Command("sc", "start", wc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Restart перезапускает сервис
func (wc *WindowsController) Restart() string {
	if wc.Stop() == "OK" {
		if wc.Start() == "OK" {
			return "OK"
		}
	}
	return "ERROR"
}

// Status возвращает статус сервиса
func (wc *WindowsController) Status() string {
	cmd := exec.Command("sc", "query", wc.leaf.Name)
	output, err := cmd.Output()
	if err != nil {
		return "UNAVAIVABLE"
	}
	
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	cnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			cnt++
			if cnt == 3 {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					status := strings.TrimSpace(parts[len(parts)-1])
					re := regexp.MustCompile(`\s+\d+\s+`)
					status = re.ReplaceAllString(status, " ")
					return strings.TrimSpace(status)
				}
			}
		}
	}
	return "UNAVAIVABLE"
}

// Logs возвращает логи сервиса (не поддерживается в Windows)
func (wc *WindowsController) Logs(args []string) string {
	return "UNAVAIVABLE"
}

// CallMethod вызывает произвольный метод
func (wc *WindowsController) CallMethod(method string, args []string) string {
	return "OK"
}
