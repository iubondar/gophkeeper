# Описание API GophKeeper

## Аутентификация и регистрация

### `POST /api/register`

Создание нового пользователя.

**Запрос:**

```json
{
  "login": "user@example.com",
  "password_hash": "<base64>",
  "salt": "<base64>"
}
```

**Ответы:**

- `200 OK` — успешно
- `409 Conflict` — пользователь уже существует
- `400 Bad Request` — неверный формат запроса

---

### `POST /api/login`

Получение соли по логину.

**Запрос:**

```json
{ "login": "user@example.com" }
```

**Ответ:**

```json
{ "salt": "<base64>" }
```

**Ответы:**

- `200 OK`
- `404 Not Found` — пользователь не найден
- `400 Bad Request`

---

### `POST /api/authenticate`

Аутентификация по логину и хэшу пароля.

**Запрос:**

```json
{
  "login": "user@example.com",
  "password_hash": "<base64>"
}
```

**Ответ:**

```json
{
  "access_token": "...",
  "refresh_token": "...",
  "expires_in": 1800
}
```

**Ответы:**

- `200 OK`
- `401 Unauthorized` — неправильный пароль
- `400 Bad Request`

---

### `POST /api/refresh`

Обновление access token с помощью refresh token. Refresh token также может быть обновлён при этом.

**Запрос:**

```json
{
  "refresh_token": "<uuid>"
}
```

**Ответ:**

```json
{
  "access_token": "...",
  "refresh_token": "...", 
  "expires_in": 1800
}
```

**Ответы:**

- `200 OK`
- `401 Unauthorized` — refresh token недействителен
- `400 Bad Request` — невалидное тело запроса

---

## Работа с записями (passwords, notes, cards)

### `POST /api/passwords`

### `POST /api/notes`

### `POST /api/cards`

Создание зашифрованной записи

> Хотя структура API одинакова для всех трёх методов, **формат содержимого ****\`\`**** (до шифрования) зависит от типа записи** и должен формироваться на клиенте.
>
> Примеры до шифрования:
>
> - Пароли: `{ "login": "user@example.com", "password": "123456" }`
> - Карты: `{ "number": "4111111111111111", "expiry": "12/26", "cvv": "123" }`
> - Заметки: `{ "title": "admin panel", "content": "https://admin.example.com" }`

**Запрос:**

```json
{
  "label": "gmail",
  "metadata": "почта",
  "data": "<base64-зашифрованный-блок>"
}
```

**Ответ:**

```json
{ "id": "abc123", "version": 1 }
```

**Ответы:**

- `200 OK`
- `409 Conflict` — label уже используется
- `400 Bad Request`
- `401 Unauthorized`

---

## Работа с файлами

### `POST /api/files`

Загрузка файла (multipart)

**Поля:**

- `label`: строка
- `metadata`: строка
- `file`: файл (зашифрованный)

**Ответ:**

```json
{ "id": "file001", "version": 1 }
```

**Ответы:**

- `200 OK`
- `409 Conflict`
- `400 Bad Request`
- `401 Unauthorized`

---

### `GET /api/files/:label/download`

Скачивание зашифрованного файла и его метаданных

**Ответ:**

- `200 OK`

  - Тело ответа: бинарный поток

- `404 Not Found`

- `401 Unauthorized`

---

## 📅 Получение и обновление записей

### `GET /api/records/:label`

Получение записи по ключу (универсально для всех типов)

**Ответ:**

```json
{
  "id": "abc123",
  "type": "password",
  "label": "gmail",
  "metadata": "...",
  "data": "<base64>",
  "version": 3
}
```

**Ответы:**

- `200 OK`
- `404 Not Found`
- `401 Unauthorized`

---

### `PUT /api/records/:label`

Обновление записи (с OCC)

**Заголовок:**

```
If-Version: 3
```

**Тело (для текстов/паролей/карт):**

```json
{
  "metadata":  "...",
  "data": "<base64>"
}
```

**Тело (для файлов, multipart):**

- `metadata`, `file`, `label` (если нужно проверить)

**Ответ:**

```json
{ "version": 4 }
```

**Ответы:**

- `200 OK`
- `409 Conflict` — версия не совпала
- `404 Not Found`
- `401 Unauthorized`

