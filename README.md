# Yggdrasil Branch - Go версия

Go версия сервиса Yggdrasil Branch для быстрого подключения сервисов к сервису Yggdrasil.

## Требования

- Go 1.21 или выше
- Доступ к systemctl (Linux) или sc (Windows) для управления сервисами
- Docker (опционально, для управления контейнерами)

## Сборка

```bash
cd go
go build -o yggdrasil-branch
```

## Запуск

```bash
# На Linux (требуются права sudo для systemctl)
sudo ./yggdrasil-branch --port=8081

# На Windows (требуются права администратора)
yggdrasil-branch.exe --port=8081
```

## Конфигурация

Сервис использует те же файлы конфигурации, что и Kotlin версия:

### websocket.config
```
# Конфигурация WebSocket подключения
websocket.host=localhost
websocket.port=8080
websocket.path=/ws
websocket.timeout=5000
```

### leaves.config
```
# Список имён systemd сервисов (Linux) или Windows сервисов
minecraft
albots
```

### docker.config
```
# Список имён docker контейнеров
gitea1
runner_1
minecraft-server
```

### user.config
```
# Глобальные права пользователей
user1:rwxa
user2:rwx
```

## API

### HTTP API

- `POST /api/leaf/connect` - Подключение листа
- `POST /api/leaf/call/{leaf}` - Вызов метода листа
- `POST /api/leaf/callback/{leaf}/{method}` - Callback от листа
- `GET /api/leaf/status` - Статус всех листов
- `GET /api/leaf/status/{leaf}` - Статус конкретного листа
- `DELETE /api/leaf/disconnect/{leaf}` - Отключение листа

### WebSocket

- `ws://localhost:8081/ws` - WebSocket endpoint для подключения листов

## Особенности

1. **Кроссплатформенность**: Работает на Linux и Windows
2. **Управление сервисами**: Поддержка systemd (Linux) и Windows Services
3. **Docker поддержка**: Управление Docker контейнерами
4. **WebSocket клиент**: Автоматическое подключение к Trunk
5. **HTTP API**: RESTful API для управления листами

## Отличия от Kotlin версии

- Использует стандартную библиотеку Go вместо Spring Boot
- Более легковесный и быстрый
- Простая сборка в один бинарный файл
- Нет зависимости от JVM

## Docker сборка

```bash
cd go
docker build -t yggdrasil-branch:latest .
```

### Запуск Docker контейнера

```bash
docker run -d \
  --name yggdrasil-branch \
  -p 8081:8080 \
  -v $(pwd)/config:/app/config \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e BRANCHNAME="your-branch-name" \
  yggdrasil-branch:latest
```

**Важно:** Для работы DockerController необходим доступ к Docker socket (`/var/run/docker.sock`).

### Docker Compose пример

```yaml
version: "3"

services:
  yggdrasil-branch:
    build: ./go
    container_name: yggdrasilBranch
    restart: always
    volumes:
      - ./config:/app/config
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      BRANCHNAME: "your-branch-name"
    ports:
      - "8081:8080"
```

## Лицензия

То же, что и основная версия проекта.
