-- +goose Up
-- +goose StatementBegin
UPDATE organizations
SET header_title = 'ДетаЛит'
WHERE header_title IS DISTINCT FROM 'ДетаЛит';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE organizations
SET header_title = 'ДетаЛит'
WHERE header_title IS NULL OR trim(header_title) = '';
-- +goose StatementEnd
