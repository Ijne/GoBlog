# 🚀 GoBlog — Modern Blogging Platform in Go

![](pictures/image1.png)

![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green?logo=postgresql)
![JWT](https://img.shields.io/badge/JWT-Auth-orange?logo=jsonwebtokens)
![Kafka](https://img.shields.io/badge/Kafka-Broker-lightgrey?logo=apachekafka)
![WS](https://img.shields.io/badge/WS-RealTime-lightgrey?logo=websocket)

**GoBlog** — это высокопроизводительный бэкенд для блоговой платформы, написанный на **Go** с использованием **PostgreSQL** и **Kafka**. Проект создан в учебных целях, но реализован с применением production-ready практик.

---

## 🔥 Особенности

✅ **REST API** для управления постами, пользователями и аутентификацией
✅ **МИКРОСЕРВИСНАЯ АРХИТЕКТУРА** логика работы с уведомлениями, вынесена в отдельный сервис
✅ **JWT-аутентификация** с защищёнными роутами
✅ **PostgreSQL** для надёжного хранения данных
✅ **Kafka** для свзяи между сервисами
✅ **Websockets** для получения уведомлений в реальном времени
✅ **Middleware** (логирование, JWT-валидация)
✅ **Конфигурация через `.env`**
✅ **Чистая архитектура** (разделение на handlers, models, storage)

---

![](pictures/image2.png)

## 🛠 Технологии

- **Язык**: Go   
- **База данных**: PostgreSQL 
- **Аутентификация**: JWT  
- **Брокер сообщений**: Kafka
- **Конфигурация**: `.env`

---

![](pictures/image3.png)

## 🚀 Быстрый старт

```bash
git clone https://github.com/Ijne/GoBlog.git
# Из GoBlog
docker compose -f 'docker-compose.yml' up -d --build


