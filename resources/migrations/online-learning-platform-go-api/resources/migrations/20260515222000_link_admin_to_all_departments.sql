-- +goose Up
-- +goose StatementBegin
INSERT INTO accounts (email, password_hash, role, username)
VALUES (
    'admin@lms.dev',
    '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6',
    'ADMIN'::roles,
    'Администратор ДетаЛит'
)
ON CONFLICT (email) DO UPDATE
SET role = 'ADMIN'::roles,
    username = EXCLUDED.username;

INSERT INTO organization_accounts (organization_id, account_id)
SELECT o.id, a.id
FROM organizations o
JOIN accounts a ON a.email = 'admin@lms.dev'
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM organization_accounts
WHERE account_id = (SELECT id FROM accounts WHERE email = 'admin@lms.dev');
-- +goose StatementEnd
