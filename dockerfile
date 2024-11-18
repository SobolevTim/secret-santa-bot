# Этап 1: Сборка приложения
#FROM golang:1.23.0-alpine AS builder
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY . .
#RUN RUN go build -o myapp ./cmd/api

# Этап 2: Создание конечного образа
#FROM alpine:latest
#WORKDIR /root/
#COPY --from=builder /app/myapp .
#COPY .env .
#CMD ["./myapp"]

# Используем официальный образ Go для сборки
FROM golang:1.23.0-alpine AS builder

# Установим рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем весь исходный код проекта в контейнер
COPY . .

# Компилируем приложение в бинарный файл
RUN go build -o myapp cmd/bot/main.go

# Используем минимальный образ на основе Alpine для уменьшения размера
FROM alpine:latest

# Установим пакет tzdata для временных зон
RUN apk add --no-cache tzdata

# Создаем рабочую директорию для контейнера
WORKDIR /root/

# Копируем скомпилированное приложение из стадии сборки
COPY --from=builder /app/myapp .

# Копируем файл .env (если используется)
COPY .env .
COPY migration/migration.sql ./migration/migration.sql

# Команда для запуска приложения
CMD ["./myapp"]