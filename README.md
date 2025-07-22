# 🚀 GoBlog — Modern Blogging Platform in Go

![](pictures/image1.png)

![Go](https://img.shields.io/badge/Go-1.22-blue?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-green?logo=postgresql)
![JWT](https://img.shields.io/badge/JWT-Auth-orange?logo=jsonwebtokens)
![GORM](https://img.shields.io/badge/GORM-ORM-lightgrey?logo=go)

**GoBlog** — это высокопроизводительный бэкенд для блоговой платформы, написанный на **Go** с использованием **PostgreSQL** и **JWT-аутентификации**. Проект создан в учебных целях, но реализован с применением production-ready практик.

---

## 🔥 Особенности

✅ **REST API** для управления постами, пользователями и аутентификацией  
✅ **JWT-аутентификация** с защищёнными роутами  
✅ **PostgreSQL** для надёжного хранения данных  
✅ **Middleware** (логирование, JWT-валидация)  
✅ **Конфигурация через `.env`**  
✅ **Чистая архитектура** (разделение на handlers, models, storage)  

---

![](pictures/image2.png)

## 🛠 Технологии

- **Язык**: Go 1.22   
- **База данных**: PostgreSQL 
- **Аутентификация**: JWT  
- **Конфигурация**: `.env`  
- **Логирование**: Встроенное + (планируется Zap)  

---

![](pictures/image3.png)

## 🚀 Быстрый старт

git clone https://github.com/Ijne/GoBlog.git
cd GoBlog/core-api_app
go run cmd/server/main.go

    ИЛИ

app.exe


