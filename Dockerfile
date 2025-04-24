# Используем официальный образ Java 21
FROM eclipse-temurin:21-jdk

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    && rm -rf /var/lib/apt/lists/*

# Установка Docker
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg \
    && echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce docker-ce-cli containerd.io \
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