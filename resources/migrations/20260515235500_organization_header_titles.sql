-- +goose Up
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS header_title varchar(80);

UPDATE organizations
SET header_title = 'ДетаЛит'
WHERE header_title IS NULL OR trim(header_title) = '';

-- +goose Down
ALTER TABLE organizations
    DROP COLUMN IF EXISTS header_title;
