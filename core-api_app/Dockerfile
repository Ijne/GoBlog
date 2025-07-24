FROM debian:bookworm-slim

WORKDIR /app
COPY app.exe /app/server
COPY .env /app/.env
EXPOSE 8080
ENTRYPOINT [ "/app/server" ]