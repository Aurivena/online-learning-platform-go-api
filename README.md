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

## Структура курсов

Система поддерживает иерархическую структуру: **Организация → Курсы → Модули → Слайды**

### API Endpoints для курсов

#### Курсы

```bash
# Список курсов организации
GET /api/organizations/:id/courses

# Создать курс
POST /api/organizations/:id/courses
{
  "title": "Название курса",
  "description": "Описание курса",
  "organization_id": 1
}

# Получить курс со всеми модулями и слайдами
GET /api/organizations/:id/courses/:courseId

# Обновить курс
PUT /api/organizations/:id/courses/:courseId
{
  "title": "Новое название",
  "description": "Новое описание"
}

# Удалить курс
DELETE /api/organizations/:id/courses/:courseId
```

#### Модули

```bash
# Создать модуль
POST /api/organizations/:id/courses/:courseId/modules
{
  "title": "Название модуля",
  "course_id": 1
}

# Получить модуль со слайдами
GET /api/organizations/:id/courses/:courseId/modules/:moduleId

# Обновить модуль
PUT /api/organizations/:id/courses/:courseId/modules/:moduleId
{
  "title": "Новое название модуля"
}

# Удалить модуль
DELETE /api/organizations/:id/courses/:courseId/modules/:moduleId

# Добавить модуль к курсу
POST /api/organizations/:id/courses/:courseId/modules
{
  "module_id": 1,
  "index": 0
}

# Удалить модуль из курса
DELETE /api/organizations/:id/courses/:courseId/modules/:moduleId
```

#### Слайды

```bash
# Создать слайд
POST /api/organizations/:id/courses/:courseId/modules/:moduleId/slides
{
  "title": "Название слайда",
  "description": "Описание",
  "slide_type": "TEXT|VIDEO_URL|TEST|FILE",
  "payload": {
    "content": "# Текст слайда",
    "videoUrl": "https://youtube.com/watch?v=...",
    "question": "Вопрос теста",
    "options": [...]
  },
  "module_id": 1
}

# Получить слайд
GET /api/organizations/:id/courses/:courseId/modules/:moduleId/slides/:slideId

# Обновить слайд
PUT /api/organizations/:id/courses/:courseId/modules/:moduleId/slides/:slideId
{
  "title": "Новое название",
  "description": "Новое описание",
  "slide_type": "TEXT",
  "payload": {...}
}

# Удалить слайд
DELETE /api/organizations/:id/courses/:courseId/modules/:moduleId/slides/:slideId

# Добавить слайд к модулю
POST /api/organizations/:id/courses/:courseId/modules/:moduleId/slides
{
  "slide_id": 1,
  "index": 0
}

# Удалить слайд из модуля
DELETE /api/organizations/:id/courses/:courseId/modules/:moduleId/slides/:slideId
```

### Типы слайдов

- **TEXT** - текстовый слайд с markdown
- **VIDEO_URL** - видео слайд с URL
- **TEST** - слайд с тестовыми вопросами
- **FILE** - слайд с файлами

### Пример создания курса с модулями

```bash
# 1. Создать модуль
POST /api/organizations/1/courses/1/modules
{
  "title": "Введение",
  "course_id": 1
}
# Response: { "id": 1, "title": "Введение" }

# 2. Создать слайд TEXT
POST /api/organizations/1/courses/1/modules/1/slides
{
  "title": "Что такое ЧПУ",
  "slide_type": "TEXT",
  "payload": {
    "content": "# ЧПУ\n\nЧисловое программное управление..."
  }
}

# 3. Создать слайд VIDEO_URL
POST /api/organizations/1/courses/1/modules/1/slides
{
  "title": "Лекция",
  "slide_type": "VIDEO_URL",
  "payload": {
    "videoUrl": "https://youtube.com/watch?v=...",
    "durationSeconds": 1200,
    "platform": "YOUTUBE"
  }
}

# 4. Получить курс со всеми модулями и слайдами
GET /api/organizations/1/courses/1
```

## Организации

Организация - это контейнер для курсов, созданная пользователем-владельцем. Каждая организация может содержать множество пользователей и курсов.

### API Endpoints для организаций

#### Управление организациями

```bash
# Создать новую организацию (текущий пользователь становится владельцем)
POST /api/organizations
{
  "title": "Название организации",
  "tag": "ORG_TAG",
  "description": "Описание организации"
}

# Список всех организаций
GET /api/organizations

# Мои организации (где я владелец)
GET /api/organizations/my

# Получить организацию по ID
GET /api/organizations/:id

# Получить организацию по тегу
GET /api/organizations/tag/:tag

# Обновить организацию
PUT /api/organizations/:id
{
  "title": "Новое название",
  "description": "Новое описание"
}

# Удалить организацию
DELETE /api/organizations/:id
```

#### Управление членами организации

```bash
# Добавить пользователя в организацию
POST /api/organizations/:id/accounts
{
  "account_id": 1
}

# Удалить пользователя из организации
DELETE /api/organizations/:id/accounts
{
  "account_id": 1
}
```

### Полный пример использования

```bash
# 1. Зарегистрироваться / Авторизоваться
POST /api/auth/register
{
  "email": "user@example.com",
  "username": "user",
  "password": "password123",
  "role": "ADMIN"
}

# 2. Создать организацию
POST /api/organizations
{
  "title": "Моя компания",
  "tag": "MYCOMPANY",
  "description": "Обучение сотрудников"
}
# Response: { "id": 1, "title": "Моя компания", ... }

# 3. Добавить членов в организацию
POST /api/organizations/1/accounts
{
  "account_id": 2
}

# 4. Создать курс
POST /api/organizations/1/courses
{
  "title": "Введение в Go",
  "description": "Учимся писать на Go",
  "organization_id": 1
}
# Response: { "id": 1, "title": "Введение в Go", ... }

# 5. Создать модуль
POST /api/organizations/1/courses/1/modules
{
  "title": "Основы Go",
  "course_id": 1
}

# 6. Создать слайд
POST /api/organizations/1/courses/1/modules/1/slides
{
  "title": "Переменные и типы",
  "slide_type": "TEXT",
  "payload": {
    "content": "# Переменные в Go\n\nВ Go переменные объявляются с помощью `var`..."
  }
}

# 7. Получить полный курс со всеми модулями и слайдами
GET /api/organizations/1/courses/1
```

### Структура данных

```
User (Account)
├── Organization (владелец)
│   ├── Members (organization_accounts)
│   └── Courses
│       ├── Module 1 (index: 1)
│       │   ├── Slide 1 (index: 1)
│       │   ├── Slide 2 (index: 2)
│       │   └── ...
│       ├── Module 2 (index: 2)
│       └── ...
```
