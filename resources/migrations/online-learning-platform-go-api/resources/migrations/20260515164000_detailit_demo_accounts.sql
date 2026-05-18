-- +goose Up
-- +goose StatementBegin
-- Пароль для всех демонстрационных учетных записей: admin
WITH demo_accounts(email, username, role) AS (
    VALUES
        ('director@detailit.ru', 'Директор обучения', 'ADMIN'::roles),
        ('master@detailit.ru', 'Мастер участка', 'ADMIN'::roles),
        ('worker@detailit.ru', 'Сотрудник производства', 'USER'::roles),
        ('admin@lms.dev', 'Администратор ДетаЛит', 'ADMIN'::roles),
        ('teacher@lms.dev', 'Мастер-наставник', 'ADMIN'::roles),
        ('student@lms.dev', 'Сотрудник производства', 'USER'::roles)
)
INSERT INTO accounts (email, username, role, password_hash, created_at, updated_at)
SELECT
    email,
    username,
    role,
    '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6',
    now(),
    now()
FROM demo_accounts
ON CONFLICT (email) DO UPDATE
SET username = EXCLUDED.username,
    role = EXCLUDED.role,
    password_hash = EXCLUDED.password_hash,
    updated_at = now();

INSERT INTO organization_accounts (organization_id, account_id)
SELECT 1, a.id
FROM accounts a
WHERE a.email IN (
    'director@detailit.ru',
    'master@detailit.ru',
    'worker@detailit.ru',
    'admin@lms.dev',
    'teacher@lms.dev',
    'student@lms.dev'
)
ON CONFLICT DO NOTHING;

SELECT setval('accounts_id_seq', (SELECT MAX(id) FROM accounts));
-- +goose StatementEnd

-- +goose Down
DELETE FROM organization_accounts
WHERE organization_id = 1
  AND account_id IN (
      SELECT id
      FROM accounts
      WHERE email IN ('admin@lms.dev', 'teacher@lms.dev', 'student@lms.dev')
  );

DELETE FROM accounts
WHERE email IN ('admin@lms.dev', 'teacher@lms.dev', 'student@lms.dev');
