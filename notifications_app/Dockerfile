FROM golang:1.25.1

WORKDIR /app

# Копируем весь проект сохраняя структуру
COPY ./ ./

RUN go mod download
RUN go build -o main cmd/server/main.go

CMD ["./main"]
