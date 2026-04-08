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
