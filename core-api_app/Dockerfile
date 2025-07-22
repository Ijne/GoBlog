FROM mcr.microsoft.com/windows/servercore:1809

WORKDIR /app
COPY app.exe /app/server
COPY .env /app/.env
EXPOSE 8080
ENTRYPOINT [ "/app/server" ]