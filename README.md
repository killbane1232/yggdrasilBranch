Yggdrasil Branch - проект для быстрого подключения сервисов к сервису Yggdrasil

Сборка:
```
./gradlew bootJar
mv ./build/libs/yggdrasilBranch-0.0.1-SNAPSHOT.jar {где удобнее будет, чтобы jar файл лежал}
```

Примеры обязательных файлов конфигурации
websocket.config:
```
# Конфигурация WebSocket подключения
# Хост для подключения WebSocket
websocket.host=localhost
# Порт для подключения WebSocket
websocket.port=8080
# Путь для WebSocket endpoint
websocket.path=/ws
# Таймаут подключения в миллисекундах
websocket.timeout=5000
```

leaves.config:
```
# Список имён windows сервисов или сервисов systemctl
minecraft
albots
```

docker.config:
```
# Список имён docker контейнеров
gitea1
runner_1
minecraft-server
```

Запуск
```
# sudo права не обязательны, если запуск идёт от пользователя, который без них может работать с systemctl
# На Windows обязательны права админинстратора
sudo java -jar ./yggdrasil.jar --server.port={порт, на котором будет работать YggdrasilBranch}
``
