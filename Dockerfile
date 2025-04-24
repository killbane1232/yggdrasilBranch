# Используем официальный образ Java 21
FROM eclipse-temurin:21-jdk

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    docker \
    && rm -rf /var/lib/apt/lists/*

# Создание рабочей директории
WORKDIR /app

# Копирование файлов проекта
COPY . .

# Сборка приложения
RUN chmod +x ./gradlew && ./gradlew bootJar
RUN mv build/libs/yggdrasilBranch-0.0.1-SNAPSHOT.jar /app/yggdrasilBranch.jar

# Создание директории для конфигурационных файлов
RUN mkdir -p /app

# Открытие порта
EXPOSE 8080

# Запуск приложения с параметром порта
CMD ["sh", "-c", "java -jar /app/yggdrasilBranch.jar --server.port=8080"] 