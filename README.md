### Makefile
- Взаимодействие с docker через Makefile
```
make dc-restart
make dc-up
```
- Взаимодействие с goose через Makefile
```
make gc_create ... -- sql файла для goose
``` 

```
make gc_up -- поднимает новые sql файлы
```

### Аутентификация и авторизация

#### Регистрация
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "role": "USER"  # или "ADMIN" - роль можно указывать с фронта (внутренний проект)
}
```

**Response:**
```json
{
  "code": 0,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "role": "USER",
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

#### Авторизация
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "code": 0,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc..."
  }
}
```

**Токены возвращаются в:**
- JSON body (access_token, refresh_token)
- HTTP cookies (httpOnly, secure)

### Файлы для запуска проекта
```.env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=my_secret_password
POSTGRES_DATABASE=lms_db
POSTGRES_HOST=localhost
POSTGRES_PORT=5433

STORAGE_USER=detalit
STORAGE_PASSWORD=secretpassword
STORAGE_PORT=9090
STORAGE_UI_PORT=9091
STORAGE_BUCKET=lms-bucket
STORAGE_ENDPOINT=localhost:9000

SERVER_ADDR=http://localhost
SERVER_PORT=8080

JWT_SECRET=super_long_secret_key_change_me

CONFIG_PATH=./resources/config.yml

ACCESS_TOKEN=...
REFRESH_TOKEN=...
```

```config.yaml
minio:
  url: ${SERVER_ADDR}:${SERVER_PORT}
  access-key: ${STORAGE_USER}
  secret-key: ${STORAGE_PASSWORD}
  bucket: ${STORAGE_BUCKET}
  endpoint: ${STORAGE_ENDPOINT}
  sslmode: false

postgres:
  user: ${POSTGRES_USER}
  password: ${POSTGRES_PASSWORD}
  database: ${POSTGRES_DATABASE}
  host: ${POSTGRES_HOST}
  port: ${POSTGRES_PORT}
  sslmode: disable

server:
  addr: ${SERVER_ADDR}
  port: ${SERVER_PORT}

token:
  access-token: ${ACCESS_TOKEN}
  refresh-token: ${REFRESH_TOKEN}
```
