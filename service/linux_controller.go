package service

import (
	"os/exec"
	"strings"
	"yggdrasil-branch/models"
)

// LinuxController управляет Linux сервисами через systemctl
type LinuxController struct {
	leaf *models.Leaf
}

// NewLinuxController создает новый LinuxController
func NewLinuxController(leaf *models.Leaf) *LinuxController {
	return &LinuxController{leaf: leaf}
}

// GetLeaf возвращает лист
func (lc *LinuxController) GetLeaf() *models.Leaf {
	return lc.leaf
}

// Stop останавливает сервис
func (lc *LinuxController) Stop() string {
	cmd := exec.Command("sudo", "/usr/bin/systemctl", "stop", lc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Start запускает сервис
func (lc *LinuxController) Start() string {
	cmd := exec.Command("sudo", "/usr/bin/systemctl", "start", lc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Restart перезапускает сервис
func (lc *LinuxController) Restart() string {
	cmd := exec.Command("sudo", "/usr/bin/systemctl", "restart", lc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Status возвращает статус сервиса
func (lc *LinuxController) Status() string {
	cmd := exec.Command("sudo", "/usr/bin/systemctl", "is-active", lc.leaf.Name)
	output, err := cmd.Output()
	if err != nil {
		return "UNAVAIVABLE"
	}
	if strings.TrimSpace(string(output)) == "active" {
		return "RUNNING"
	}
	return "STOPPED"
}

// Logs возвращает логи сервиса
func (lc *LinuxController) Logs(args []string) string {
	num := 10
	if len(args) > 0 {
		if n, err := parseInt(args[0]); err == nil && n > 0 {
			num = n
		}
	}
	
	cmd := exec.Command("sudo", "/usr/bin/journalctl", "--unit="+lc.leaf.Name, "-n", intToString(num), "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return "UNAVAIVABLE"
	}
	return string(output)
}

// CallMethod вызывает произвольный метод
func (lc *LinuxController) CallMethod(method string, args []string) string {
	return "OK"
}
