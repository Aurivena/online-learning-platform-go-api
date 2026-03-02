CREATE TYPE roles as ENUM ('USER','ADMIN');

CREATE TABLE accounts
(
    id            BIGSERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(50)  NOT NULL,
    username      VARCHAR(125),
    login         VARCHAR(125),
    created_at    timestamp default now()
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

CREATE TABLE organizations
(
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(125) NOT NULL,
    tag         VARCHAR(15)  NOT NULL UNIQUE,
    description text         NOT NULL,
    owner_id    BIGINT       NOT NULL,
    created_at  timestamp default now()
);


CREATE TABLE organization_accounts
(
    organization_id BIGINT,
    account_id      BIGINT,

    CONSTRAINT fk_organization_accounts_0
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_organization_accounts_1
        FOREIGN KEY (account_id)
            REFERENCES accounts (id)
            ON DELETE CASCADE,
    PRIMARY KEY (organization_id, account_id)
);

-- (пароль 'admin' в BCrypt)
INSERT INTO accounts (email, password_hash, role, username, login)
VALUES ('admin@lms.dev', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'ADMIN', 'Admin',
        'admin_boss');

INSERT INTO organizations (title, tag, description, owner_id, created_at)
VALUES ('org1-test','HEROBRIN','Майн',1,now());

INSERT INTO organization_accounts (organization_id, account_id)
VALUES (1,1);
