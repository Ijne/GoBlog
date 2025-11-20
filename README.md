# 🚀 GoBlog — Modern Blogging Platform in Go

![](pictures/image1.png)

![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green?logo=postgresql)
![JWT](https://img.shields.io/badge/JWT-Auth-orange?logo=jsonwebtokens)
![Kafka](https://img.shields.io/badge/Kafka-Broker-lightgrey?logo=apachekafka)
![WS](https://img.shields.io/badge/WS-RealTime-lightgrey?logo=websocket)

**GoBlog** is a high-performance backend for a blogging platform built in **Go** using **PostgreSQL** and **Kafka**. The project was created for educational purposes but implements production-ready practices.

---

## 🔥 Features

✅ **REST API** for managing posts, users, and authentication  
✅ **MICROSERVICE ARCHITECTURE** - notification logic is separated into a dedicated service  
✅ **JWT authentication** with protected routes  
✅ **PostgreSQL** for reliable data storage  
✅ **Kafka** for inter-service communication  
✅ **WebSockets** for real-time notifications  
✅ **Middleware** (logging, JWT validation)  
✅ **Configuration via `.env`**  
✅ **Clean architecture** (separation into handlers, models, storage)  

---

![](pictures/image2.png)

## 🛠 Technologies

- **Language**: Go  
- **Database**: PostgreSQL  
- **Authentication**: JWT  
- **Message Broker**: Kafka  
- **Configuration**: `.env`  

---

![](pictures/image3.png)

## 🚀 Quick Start

```bash
git clone https://github.com/Ijne/GoBlog.git
# From GoBlog directory
docker compose -f 'docker-compose.yml' up -d --build
<<<<<<< HEAD
# First launch may take 1-2 minutes, subsequent launches ~30 seconds
```
=======
# Первый запуск может занять 1-2 минуты, затем ~30 секунд


>>>>>>> 6b9ba11 (Update README)
