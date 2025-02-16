# Магазин мерча

## API Эндпоинты

### 1. Авторизация

**POST /api/auth**  
Аутентификация пользователя и получение токена доступа.

#### Запрос:

```json
{
  "username": "user1",
  "password": "user1"
}
```

#### Ответ:

```json
{
  "response": {
    "accessToken": "your_jwt_token"
  }
}
```

#### Возможные ошибки:

- 400 Bad Request - если переданы некорректные данные
- 401 Unauthorized - если неверный пароль
- 500 Internal Server Error - ошибка сервера

### 2. Получение информации о переводах и покупках пользователя

**GET /api/info**  
Получение истории покупок и транзакций пользователя.

#### cookie:

```json
Cookie: accessToken=your_jwt_token
```

#### Ответ:

```json
{
  "purchases": [
    {
      "guid": "f3d140c4-30a8-451a-b578-0d837f9d9300",
      "user_id": "user1",
      "merch_name": "t-shirt",
      "created_at": "2025-02-16T20:38:53.414706+03:00"
    }
  ],
  "transactions": [
    {
      "guid": "",
      "created_at": "2025-02-16T20:38:53.42056+03:00",
      "receiver_username": "user1",
      "sender_username": "user2",
      "money_amount": 100
    }
  ]
}
```

#### Возможные ошибки:

- 400 Bad Request - если отсутствует токен или данные повреждены
- 500 Internal Server Error - ошибка сервера
- 401 Unauthorized - если не авторизован


### 3. Отправка монет другому пользователю

**POST /api/sendCoin**  
Отправка определённого количества монет другому пользователю

#### cookie:

```json
Cookie: accessToken=your_jwt_token
```

#### Запрос:

```json
{
  "receiver_username": "user2",
  "amount": 30.0
}
```

#### Ответ:

```json
{
  "guid": "768301ec-de59-4a90-a500-de477e7c7746",
  "created_at": "2025-02-16T19:16:35.999232+03:00",
  "receiver_username": "user2",
  "sender_username": "user1",
  "money_amount": 30
}
```

#### Возможные ошибки:

- 400 Bad Request - если переданы некорректные данные или недостаточно средств.
- 401 Unauthorized - если не авторизован
- 500 Internal Server Error - ошибка сервера.



### 4.Покупка товара

**POST /api/buy/:item**  
Позволяет пользователю купить товар, списывая соответствующую сумму с баланса.

#### cookie:

```json
Cookie: accessToken=your_jwt_token
```

#### Ответ:

```json
{
  "guid": "91f26f19-6ba3-4203-a851-021913cec6a8",
  "user_id": "user1",
  "merch_name": "socks",
  "created_at": "2025-02-16T16:39:17.662729803Z"
}
```

#### Возможные ошибки:

- 400 Bad Request - если переданы некорректные данные или недостаточно средств.
- 401 Unauthorized - если не авторизован
- 500 Internal Server Error - ошибка сервера.


# Авторизация
Используется JWT-токен, который передаётся в Cookie при каждом запросе. Срок действия токена — 5 часов.

## Дополнительно 
- Для оптимизациии запросов были использованы индексы
- Чтобы предотвратить грязное чтение, используеются транзакции при переводе coins и при покупке мерча
- Настроен ci на запуск тестов и линтера при push и pull request в master

# Тесты

- unit тесты для каждого слоя
- интеграционные тесты для :
    - покупки мерча
    - передачи coins другим сотрудникам
    - авторизации
    - получения информации о покупках и транзакциях
- нагрузочное тестирование
- покрытие кода тестами можно посмотреть в coverage.html

# Запуск в Docker

Склонировать проект с гита

```
git clone https://github.com/aleksandra0KR/shop
```

Перейти в директорию проекта

```
cd shop
```

Забилдить

```
docker compose build
```

Запустить:

```
docker compose up
```

---

# Запустить без Docker

```
git clone https://github.com/aleksandra0KR/shop
```

Перейти в директорию проекта

```
cd shop
```

Запустить

```
go run cmd/main.go
```

### В файле .env можно поменять на нужные вам параметры