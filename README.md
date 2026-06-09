<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=0:0B6B6B,100:0D4F4F&height=200&section=header&text=NP&fontSize=80&fontColor=ffffff&animation=fadeIn" alt="NP" />
</p>

<h3 align="center">
  Mauritius Wildlife — E-Commerce & Donation Platform
</h3>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white" alt="Go 1.26"/>
  <img src="https://img.shields.io/badge/Rust-1.85-dea584?logo=rust&logoColor=white" alt="Rust"/>
  <img src="https://img.shields.io/badge/Axum-0.8-8B5CF6" alt="Axum"/>
  <img src="https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL"/>
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" alt="Docker Compose"/>
  <img src="https://img.shields.io/badge/JWT-auth-000000?logo=jsonwebtokens" alt="JWT Auth"/>
  <img src="https://img.shields.io/badge/license-MIT-yellow" alt="License"/>
</p>

<p align="center">
  <b>NP</b> — полнофункциональная платформа электронной коммерции и пожертвований,<br/>
  посвящённая морской и дикой природе Маврикия.
</p>

---

## Overview

NP — это полиглотный микросервисный проект, объединяющий **Go** и **Rust** в единую экосистему. Платформа включает каталог товаров, корзину покупок, систему платежей, кошелёк для пожертвований и историю транзакций.

Все сервисы оркестрируются через Docker Compose, аутентификация построена на JWT (HS256) с хранением в cookie, база данных — PostgreSQL с GORM AutoMigrate.

---

## Architecture

```
                    ┌─────────────┐
                    │  API Gateway │  (Go — :8080)
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          ▼                ▼                ▼
   ┌────────────┐   ┌────────────┐   ┌────────────┐
   │ Web Service │   │User Service│   │Bank Service │
   │ (Rust/Axum) │   │  (Go)      │   │  (Go)       │
   │   :8081     │   │  :8084     │   │  :8083      │
   └────────────┘   └────────────┘   └────────────┘
                                       │
                              ┌────────┴────────┐
                              │  Order Service   │
                              │     (Go)         │
                              │     :8082        │
                              └─────────────────┘
                                       │
                              ┌────────┴────────┐
                              │   PostgreSQL     │
                              │     :5432        │
                              └─────────────────┘
```

### Services

| Service | Language | Port | Description |
|---------|----------|------|-------------|
| `api-gateway` | Go | 8080 | Reverse-proxy gateway, JWT middleware |
| `web-service` | Rust (Axum) | 8081 | Frontend pages (index, about, buy-merch) |
| `order-service` | Go | 8082 | Cart management, purchases |
| `bank-service` | Go | 8083 | Wallet, donations, top-up, transactions |
| `user-service` | Go | 8084 | Registration & authentication |

---

## Tech Stack

**Backend:**
- [Go](https://go.dev/) 1.26 — микросервисы, `gorilla/mux`, `gorm.io`
- [Rust](https://www.rust-lang.org/) edition 2024 — web-service на `axum 0.8`
- [PostgreSQL](https://www.postgresql.org/) 15 — основная БД
- [Docker Compose](https://docs.docker.com/compose/) — оркестрация

**Authentication & Security:**
- JWT (HS256, 24h expiry, cookie-based)
- bcrypt password hashing

**Frontend:**
- Server-side rendered Go templates
- Ocean/teal design system (inline CSS)

---

## Features

- **🛍️ Каталог товаров** — просмотр и покупка мерча
- **🛒 Корзина** — добавление, удаление, оформление
- **💰 Кошелёк** — пополнение баланса, история операций
- **🎁 Донаты** — пожертвования через банковский сервис
- **🔐 JWT-аутентификация** — безопасная регистрация и вход
- **📱 Адаптивный UI** — тёмно-бирюзовая тема океана

---

## Quick Start

```bash
# 1. Клонирование
git clone https://github.com/kempedron/NP.git
cd NP

# 2. Запуск
./start.sh
```

> `start.sh` выполняет: `docker-compose down && docker-compose build --parallel && docker-compose up`

После запуска сервис доступен по адресу: [http://localhost:8080](http://localhost:8080)

---

## Project Structure

```
NP/
├── cmd/
│   ├── api-gateway/       # Go — точка входа API Gateway
│   ├── web-service/       # Rust — frontend сервис
│   ├── order-service/     # Go — управление заказами
│   ├── bank-service/      # Go — банковский сервис
│   └── user-service/      # Go — пользовательский сервис
├── internal/
│   ├── analytic-service/  # stubs для аналитики
│   ├── database/          # GORM инициализация + миграции
│   ├── jwt/               # JWT генерация/валидация
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Модели данных (GORM)
│   └── ...                # Бизнес-логика по сервисам
├── web/
│   ├── static/            # Статические файлы
│   └── templates/         # HTML шаблоны (Go + Rust)
├── docker-compose.yml     # Оркестрация сервисов
└── start.sh               # Скрипт запуска
```

---

## API Endpoints

| Method | Path | Service | Description |
|--------|------|---------|-------------|
| `GET` | `/` | web-service | Главная страница |
| `GET` | `/about-us` | web-service | О нас |
| `GET` | `/buy-merch` | web-service | Каталог товаров |
| `GET/POST` | `/register` | user-service | Регистрация |
| `GET/POST` | `/login` | user-service | Вход |
| `GET` | `/cart` | order-service | Корзина |
| `POST` | `/cart/add` | order-service | Добавить в корзину |
| `POST` | `/cart/delete` | order-service | Удалить из корзины |
| `POST` | `/cart/purchase` | order-service | Оформление покупки |
| `GET` | `/my-purchases` | order-service | История покупок |
| `GET` | `/wallet` | bank-service | Кошелёк |
| `POST` | `/wallet/top-up` | bank-service | Пополнение баланса |
| `POST` | `/wallet/donate` | bank-service | Донат |

---


## License

Распространяется под лицензией **MIT**. См. файл [LICENSE](LICENSE) для получения дополнительной информации.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/kempedron">kempedron</a>
</p>

<p align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=0:0D4F4F,100:0B6B6B&height=120&section=footer" />
</p>
