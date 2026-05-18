-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS course_organizations
(
    course_id       BIGINT NOT NULL,
    organization_id BIGINT NOT NULL,
    created_at      timestamp default now(),

    CONSTRAINT fk_course_organizations_course_id
        FOREIGN KEY (course_id)
            REFERENCES courses (id) ON DELETE CASCADE,

    CONSTRAINT fk_course_organizations_organization_id
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id) ON DELETE CASCADE,

    PRIMARY KEY (course_id, organization_id)
);

CREATE INDEX IF NOT EXISTS idx_course_organizations_organization_id ON course_organizations (organization_id);
CREATE INDEX IF NOT EXISTS idx_course_organizations_course_id ON course_organizations (course_id);

INSERT INTO course_organizations (course_id, organization_id)
SELECT id, organization_id
FROM courses
WHERE organization_id IS NOT NULL
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_organizations;
-- +goose StatementEnd
