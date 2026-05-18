-- +goose Up
-- +goose StatementBegin
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
    organization_id BIGINT NOT NULL,
    account_id      BIGINT NOT NULL,

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

CREATE INDEX idx_organization_accounts_account_id ON organization_accounts (account_id);
CREATE INDEX idx_organization_accounts_organization_id ON organization_accounts (organization_id);

INSERT INTO organizations (id, title, tag, description, owner_id, created_at)
VALUES (1,
        'ООО "ДетаЛит"',
        'detailit',
        'Производственное подразделение для обучения сотрудников литейного участка, механической обработки и отдела технического контроля.',
        1,
        now());

INSERT INTO organization_accounts (organization_id, account_id)
VALUES (1, 1),
       (1, 2),
       (1, 3);

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS organization_accounts;
DROP TABLE IF EXISTS organizations;
