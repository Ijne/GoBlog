FROM golang:1.25.1

WORKDIR /app

COPY ./notifications_app .

CMD ["go", "run", "notifications_app/cmd/server/main.go"]
