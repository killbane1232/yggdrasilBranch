# Стадия сборки
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /build

# Копирование go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o yggdrasil-branch .

# Стадия выполнения
FROM alpine:3.21

# Установка необходимых пакетов
RUN apk add --no-cache \
    ca-certificates \
    docker \
    tzdata \
    && rm -rf /var/cache/apk/*

# Создание пользователя для запуска приложения (не root)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Создание директории для приложения
WORKDIR /app

# Копирование бинарного файла из стадии сборки
COPY --from=builder /build/yggdrasil-branch /app/yggdrasil-branch

# Создание директории для конфигурационных файлов
RUN mkdir -p /app/config

# Изменение владельца файлов
RUN chown -R appuser:appuser /app

# Переключение на непривилегированного пользователя
USER appuser

# Открытие порта
EXPOSE 8080

# Переменные окружения
ENV PORT=8080

# Запуск приложения
CMD ["./yggdrasil-branch", "--port=8080"]
