-- +goose Up
-- +goose StatementBegin
UPDATE slides
SET payload = jsonb_set(
    payload,
    '{videoUrl}',
    to_jsonb('https://www.youtube.com/watch?v=V9RLc9EX1so'::text),
    true
)
WHERE slide_type = 'VIDEO_URL'
  AND payload ->> 'videoUrl' = 'https://www.youtube.com/watch?v=g79igk3edru';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE slides
SET payload = jsonb_set(
    payload,
    '{videoUrl}',
    to_jsonb('https://www.youtube.com/watch?v=g79igk3edru'::text),
    true
)
WHERE slide_type = 'VIDEO_URL'
  AND payload ->> 'videoUrl' = 'https://www.youtube.com/watch?v=V9RLc9EX1so';
-- +goose StatementEnd
