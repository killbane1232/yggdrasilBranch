# HTTP API для LeafHttpController

`LeafHttpController` предоставляет HTTP API для работы с листами, аналогично WebSocket функциональности в `LeafController`.

## Базовый URL
```
http://localhost:8080/api/leaf
```

## Эндпоинты

### 1. Подключение листа
**POST** `/api/leaf/connect`

Подключает лист к ветке.

**Тело запроса:**
```json
{
  "name": "leaf-name",
  "status": "active",
  "hooks": []
}
```

**Ответ:**
```json
{
  "name": "leaf-name",
  "status": "active",
  "attachedBranch": "branch-name",
  "hooks": []
}
```

### 2. Вызов метода листа
**POST** `/api/leaf/call/{leaf}`

Вызывает метод на указанном листе.

**Параметры пути:**
- `leaf` - имя листа

**Тело запроса:**
```json
{
  "method": "methodName",
  "args": ["arg1", "arg2"]
}
```

**Ответ:** Результат выполнения метода в виде строки.

### 3. Callback от листа
**POST** `/api/leaf/callback/{leaf}/{method}`

Обрабатывает callback от листа.

**Параметры пути:**
- `leaf` - имя листа
- `method` - имя метода

**Тело запроса:** Строка с ответом

**Ответ:** `"Callback processed"`

### 4. Получение статуса всех листов
**GET** `/api/leaf/status`

Возвращает список всех подключенных листов.

**Ответ:**
```json
[
  {
    "name": "leaf1",
    "status": "active",
    "attachedBranch": "branch-name",
    "hooks": []
  },
  {
    "name": "leaf2",
    "status": "inactive",
    "attachedBranch": "branch-name",
    "hooks": []
  }
]
```

### 5. Получение статуса конкретного листа
**GET** `/api/leaf/status/{leaf}`

Возвращает статус конкретного листа.

**Параметры пути:**
- `leaf` - имя листа

**Ответ:**
```json
{
  "name": "leaf-name",
  "status": "active",
  "attachedBranch": "branch-name",
  "hooks": []
}
```

**Ошибка 404:** Если лист не найден.

### 6. Отключение листа
**DELETE** `/api/leaf/disconnect/{leaf}`

Отключает указанный лист.

**Параметры пути:**
- `leaf` - имя листа

**Ответ:** `"Leaf {leaf} disconnected"`

## Примеры использования

### Подключение листа
```bash
curl -X POST http://localhost:8080/api/leaf/connect \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-leaf",
    "status": "active",
    "hooks": []
  }'
```

### Вызов метода
```bash
curl -X POST http://localhost:8080/api/leaf/call/test-leaf \
  -H "Content-Type: application/json" \
  -d '{
    "method": "getStatus",
    "args": []
  }'
```

### Получение статуса
```bash
curl http://localhost:8080/api/leaf/status
```

## Особенности

1. **Совместимость с WebSocket:** HTTP API использует тот же `RequestBuffer` и `LeafCollector`, что обеспечивает совместимость с существующей WebSocket функциональностью.

2. **Синхронные вызовы:** В отличие от WebSocket, HTTP API предоставляет синхронные вызовы методов.

3. **Автоматическая очистка:** Контроллер автоматически очищает устаревшие запросы каждую минуту.

4. **Потокобезопасность:** Все операции с `leafCollector` синхронизированы для обеспечения потокобезопасности.

