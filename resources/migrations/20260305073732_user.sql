-- +goose Up
-- +goose StatementBegin
CREATE TYPE roles as ENUM ('USER','ADMIN');

CREATE TABLE accounts
(
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          roles  default 'USER',
    username      VARCHAR(125),
    created_at    timestamp default now(),
    updated_at    timestamp default now()
);

CREATE TABLE refresh_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    token_hash VARCHAR(255) NOT NULL,
    account_id BIGINT       NOT NULL,
    expiration timestamp    NOT NULL,
    created    timestamp    NOT NULL,

    CONSTRAINT fk_refresh_token_to_account_0
        FOREIGN KEY (account_id)
            REFERENCES accounts (id)
            ON DELETE CASCADE
);

-- Все демонстрационные пользователи используют пароль 'admin'.
INSERT INTO accounts (id, email, password_hash, role, username)
VALUES (1, 'director@detailit.ru', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'ADMIN', 'Директор обучения'),
       (2, 'master@detailit.ru',   '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'ADMIN', 'Мастер участка'),
       (3, 'worker@detailit.ru',   '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'USER',  'Сотрудник производства');

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS accounts;
DROP TYPE IF EXISTS roles;
