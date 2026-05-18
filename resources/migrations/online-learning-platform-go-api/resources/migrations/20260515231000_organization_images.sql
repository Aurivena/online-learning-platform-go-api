-- +goose Up
-- +goose StatementBegin
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS image_url text;

UPDATE organizations
SET image_url = CASE
    WHEN tag IN ('detailit', 'det-prod', 'det-line1', 'det-line2', 'det-proto')
        THEN '/api/files/course_files/detailit/electronics-architecture.png'
    WHEN tag IN ('det-harness', 'det-log', 'det-plan')
        THEN '/api/files/course_files/detailit/harness-assembly.png'
    WHEN tag IN ('det-safety', 'det-eng', 'det-smt', 'det-pcb', 'det-store')
        THEN '/api/files/course_files/detailit/esd-workstation.png'
    WHEN tag IN ('det-lab', 'det-qc', 'det-calib', 'det-repair', 'det-pack', 'det-maint', 'det-train')
        THEN '/api/files/course_files/detailit/testing-traceability.png'
    ELSE '/api/files/course_files/detailit/process-flow.png'
END
WHERE image_url IS NULL OR image_url = '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE organizations
    DROP COLUMN IF EXISTS image_url;
-- +goose StatementEnd
