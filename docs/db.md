# Структура базы данных GophKeeper

## Таблица `users`

Хранит информацию о пользователях и данных, необходимых для верификации пароля.

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  login TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  salt TEXT NOT NULL
);
```

| Поле            | Тип  | Описание                        |
| --------------- | ---- | ------------------------------- |
| `id`            | UUID | Уникальный идентификатор        |
| `login`         | TEXT | Логин пользователя (уникальный) |
| `password_hash` | TEXT | Хэш пароля (после KDF)          |
| `salt`          | TEXT | Соль для генерации ключей       |

---

## Таблица `records`

Хранит все записи пользователя: пароли, текст, карты и файлы (в виде ссылок).

```sql
CREATE TABLE records (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  label TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL,
  metadata TEXT,
  encrypted_data BYTEA,
  file_key TEXT,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
```

| Поле             | Тип         | Описание                                        |
| ---------------- | ----------- | ----------------------------------------------- |
| `id`             | UUID        | Уникальный идентификатор записи                 |
| `user_id`        | UUID        | Владелец записи                                 |
| `label`          | TEXT        | Ключ доступа (уникальный в рамках пользователя) |
| `type`           | TEXT        | Тип: `password`, `text`, `card`, `file`         |
| `metadata`       | TEXT        | Метаданные в plain text                         |
| `encrypted_data` | BYTEA       | Зашифрованный blob (не используется для файлов) |
| `file_key`       | TEXT        | Ключ файла в MinIO (если `type = file`)         |
| `version`        | INTEGER     | Версия записи (для OCC)                         |
| `created_at`     | TIMESTAMPTZ | Дата создания                                   |
| `updated_at`     | TIMESTAMPTZ | Дата последнего обновления                      |

---

## 🔍 Индексы

```sql
CREATE INDEX idx_records_user_label ON records(user_id, label);
```



