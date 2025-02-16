# Avito Shop Service

## Введение

Этот сервис представляет собой внутреннюю систему обмена монетами и покупки мерча внутри компании.\
Сотрудники могут отправлять друг другу монеты, покупать мерч и отслеживать историю своих транзакций.\
Были выполнены все основные и дополнительные [условия задания](https://github.com/avito-tech/tech-internship/blob/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/Backend-trainee-assignment-winter-2025.md),
за исключением нагрузочного тестирования. В дополнении, было реализовано покрытие тестами слоев, отвечающих за работу с входящими запросами и работу с базой данных. 
При написании проекта придерживался приемов чистой архитектуры. 

### Основные возможности:

- **Аутентификация JWT-токен**
- **Отправка монет другим пользователям**
- **Покупка мерча за монеты**
- **Просмотр баланса, инвентаря и истории транзакций**

### Используемые технологии:

- **Язык:** Golang 1.23
- **База данных:** PostgreSQL
- **Фреймворки и библиотеки:**
    - [gin-gonic](https://github.com/gin-gonic/gin) – роутер и обработка запросов
    - [viper](https://github.com/spf13/viper) – конфигурация
    - [sqlx](https://github.com/jmoiron/sqlx) – работа с БД
    - [go-transaction-manager](https://github.com/avito-tech/go-transaction-manager) - работа с транзакциями
    - [jwt-go](https://github.com/dgrijalva/jwt-go) – авторизация
    - [logrus](https://github.com/sirupsen/logrus) – логирование
- **Контейнеризация:** Docker Compose
- **Тестирование:** [testify](https://github.com/stretchr/testify), [gomock](https://github.com/golang/mock), [sqlmock](https://github.com/DATA-DOG/go-sqlmock)

---

## Начало работы

### Установка

1. **Клонировать репозиторий:**

   ```sh
   git clone https://github.com/senyabanana/avito-shop-service.git
   cd avito-shop-service
   ```
   
2. **Создать и заполнить файл `.env` на основе `.env.example`**


3. **Запустить сервис через Docker Compose:**

   ```sh
   docker-compose up --build -d
   ```

### Использование

| **Команда**                        | **Описание**                 |
|------------------------------------| ---------------------------- |
| `docker-compose up --build -d`     | Запустить сервис             |
| `docker-compose down`              | Остановить сервис            |
| `docker-compose down -v`           | Остановить и очистить volume |

### Тестирование

1. **Unit-тестирование:**

   ```sh
   go test -v ./internal/... --cover
   ```
   
    Для дополнительной информации о покрытии, выполните:

   ```sh
   go test -coverprofile=coverage.out ./internal/... && go tool cover -html=coverage.out
   ```

2. **E2E-тестирование:**

    Для данного тестирования, убедитесь, что сервис запущен и выполните:

   ```sh
   go test -v -tags=e2e ./tests/e2e
   ```

---

## Линтинг кода

Проект использует `golangci-lint` для анализа кода, предотвращения ошибок и обеспечения единообразия кодовой базы.

### Конфигурация `.golangci.yml`

#### **Линтеры:**

- **govet** – стандартный анализатор кода Go, выявляющий потенциальные ошибки.
- **staticcheck** – расширенный анализатор, находящий устаревший код и потенциальные ошибки.
- **errcheck** – проверяет, что все ошибки обработаны корректно.
- **gosimple** – упрощает код, предлагая более идиоматичные конструкции.
- **ineffassign** – обнаруживает переменные, которым присвоены, но не использованные значения.
- **unused** – проверяет неиспользуемые переменные, константы и типы.
- **gofmt** – автоматически форматирует код в соответствии со стандартами Go.
- **revive** – гибкая замена `golint`, позволяющая настраивать правила проверки кода.
- **gci** – проверяет порядок импортов и автоматически группирует их.
- **bodyclose** – проверяет закрытие тел HTTP-запросов.
- **dogsled** – проверяет пропущенные идентификаторы в множественных присваиваниях.

#### **Настройки линтеров:**

- `gci` – определяет порядок импортов (стандартные, сторонние, локальные).
- `revive` – содержит правила именования переменных, комментариев к пакетам и предупреждения о проблемах в коде.

#### **Запуск линтера:**

   ```sh
     golangci-lint run
   ```

---

## Эндпоинты API

### **Аутентификация**

#### `POST /api/auth`

- **Описание:** Вход пользователя. При первой аутентификации аккаунт создается автоматически.
- **Тело запроса:**
  ```json
  {
    "username": "user1",
    "password": "password123"
  }
  ```
- **Тело ответа (успех 200 OK):**
  ```json
  {
    "token": "jwt-token"
  }
  ```
- **Ошибки:**
    - `400 Bad Request` – Неверный формат запроса
    - `401 Unauthorized` – Ошибка авторизации
    - `500 Internal Server Error` – Ошибка сервера

---

### **Получение информации**

#### `GET /api/info`

- **Описание:** Возвращает баланс пользователя, инвентарь и историю транзакций.
- **Требуется Bearer-токен в заголовке.**
- **Тело ответа (успех 200 OK):**
  ```json
  {
    "coins": 1000,
    "inventory": [
      {
        "type": "t-shirt",
        "quantity": 1
      }
    ],
    "coinHistory": {
      "received": [
        {
          "fromUser": "alice",
          "amount": 50
        }
      ],
      "sent": [
        {
          "toUser": "bob",
          "amount": 20
        }
      ]
    }
  }
  ```
- **Ошибки:**
    - `401 Unauthorized` – Токен отсутствует или невалиден
    - `500 Internal Server Error` – Ошибка сервера

---

### **Отправка монет**

#### `POST /api/sendCoin`

- **Описание:** Отправить монеты другому пользователю.
- **Тело запроса:**
  ```json
  {
    "toUser": "bob",
    "amount": 100
  }
  ```
- **Тело ответа (успех 200 OK):**
  ```json
  {
    "status": "coins were successfully sent to the user"
  }
  ```
- **Ошибки:**
    - `400 Bad Request` – Некорректные данные (некорректный пользователь, недостаточно монет)
    - `401 Unauthorized` – Ошибка авторизации
    - `500 Internal Server Error` – Ошибка сервера

---

### **Покупка мерча**

#### `GET /api/buy/{item}`

- **Описание:** Покупка мерча за монеты.
- **Пример запроса:** `/api/buy/t-shirt`
- **Тело ответа (успех 200 OK):**
  ```json
  {
    "status": "item was successfully purchased"
  }
  ```
- **Ошибки:**
    - `400 Bad Request` – Некорректные данные (товар не найден, недостаточно монет)
    - `401 Unauthorized` – Ошибка авторизации
    - `500 Internal Server Error` – Ошибка сервера

---

## Проблемы, с которыми столкнулся

- **Коды ошибок в эндпоинте `/api/info`**

    - Исключил `400 Bad Request`, так как он не актуален для этого эндпоинта.

- **Тело ответа в эндпоинтах `/api/sendCoin` и `/api/buy/{item}`**

    - Реализовал тело ответа при успешном выполнении запроса, чтобы было понимание, что запрос выполнился без ошибок.

---
