package service

import (
	"bufio"
	"os/exec"
	"strings"
	"yggdrasil-branch/models"
)

// DockerController управляет Docker контейнерами
type DockerController struct {
	leaf *models.Leaf
}

// NewDockerController создает новый DockerController
func NewDockerController(leaf *models.Leaf) *DockerController {
	return &DockerController{leaf: leaf}
}

// GetLeaf возвращает лист
func (dc *DockerController) GetLeaf() *models.Leaf {
	return dc.leaf
}

// Stop останавливает контейнер
func (dc *DockerController) Stop() string {
	cmd := exec.Command("docker", "stop", dc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Start запускает контейнер
func (dc *DockerController) Start() string {
	cmd := exec.Command("docker", "start", dc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Restart перезапускает контейнер
func (dc *DockerController) Restart() string {
	cmd := exec.Command("docker", "restart", dc.leaf.Name)
	if err := cmd.Run(); err != nil {
		return "ERROR"
	}
	return "OK"
}

// Status возвращает статус контейнера
func (dc *DockerController) Status() string {
	cmd := exec.Command("docker", "ps", "-f", "name="+dc.leaf.Name)
	output, err := cmd.Output()
	if err != nil {
		return "UNAVAIVABLE"
	}
	
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Up") {
			return "RUNNING"
		}
	}
	return "STOPPED"
}

// Logs возвращает логи контейнера
func (dc *DockerController) Logs(args []string) string {
	num := 10
	if len(args) > 0 {
		if n, err := parseInt(args[0]); err == nil && n > 0 {
			num = n
		}
	}
	
	cmd := exec.Command("docker", "logs", dc.leaf.Name, "-n", intToString(num))
	output, err := cmd.Output()
	if err != nil {
		return "UNAVAIVABLE"
	}
	return string(output)
}

// CallMethod вызывает произвольный метод
func (dc *DockerController) CallMethod(method string, args []string) string {
	return "OK"
}
