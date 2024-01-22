# Используем официальный образ Golang
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go модули и файлы проекта
COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта
COPY . .

# Собираем приложение
RUN go build -o my-go-app ./cmd
RUN chmod +x my-go-app   # Устанавливаем права на выполнение

# Команда для запуска приложения
CMD ["./my-go-app"]
