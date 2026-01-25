# =========================
# Стадия сборки
# =========================
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


# =========================
# Стадия выполнения
# =========================
FROM alpine:3.21

# 👇 GID группы docker на хосте (по умолчанию)
ARG DOCKER_GID=998

RUN apk add --no-cache \
    ca-certificates \
    docker-cli \
    tzdata \
    shadow \
 && addgroup -g ${DOCKER_GID} docker \
 && addgroup -g 1000 appuser \
 && adduser -D -u 1000 -G appuser appuser \
 && adduser appuser docker \
 && rm -rf /var/cache/apk/*

# Создание пользователя для запуска приложения (не root)
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Создание директории для приложения
WORKDIR /app

# Копирование бинарного файла из стадии сборки
COPY --from=builder /build/yggdrasil-branch /app/yggdrasil-branch

RUN mkdir -p /app/config \
 && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# Переменные окружения
ENV PORT=8080

CMD ["./yggdrasil-branch", "--port=8080"]
