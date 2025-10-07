FROM golang:1.25.1

WORKDIR /app

COPY . .

RUN ls

RUN go build -o main cmd/server/main.go

RUN ls

CMD ["./main"]
