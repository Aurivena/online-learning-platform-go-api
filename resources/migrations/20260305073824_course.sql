-- +goose Up
-- +goose StatementBegin
CREATE TYPE slide_variation as ENUM ('TEXT','VIDEO_URL','TEST','FILE');

CREATE TABLE courses
(
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     text         NOT NULL,
    owner           BIGINT,
    organization_id BIGINT,
    created_at      timestamp default now(),

    CONSTRAINT fk_course_owner_id
        FOREIGN KEY (owner)
            REFERENCES accounts (id) ON DELETE CASCADE ,

    CONSTRAINT fo_course_organization_id
        FOREIGN KEY (organization_id)
            REFERENCES organizations (id) ON DELETE CASCADE
);

CREATE TABLE modules
(
    id         BIGSERIAL PRIMARY KEY,
    title      VARCHAR(125) NOT NULL,
    created_at timestamp default now()
);

CREATE TABLE slides
(
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(255),
    description text,
    slide_type  slide_variation NOT NULL,
    payload     jsonb,
    created_at  timestamp default now()
);

CREATE TABLE course_modules
(
    course_id BIGINT,
    module_id BIGINT,
    index     int,

    CONSTRAINT fk_course_modules_course_id
        FOREIGN KEY (course_id)
            REFERENCES courses (id) ON DELETE CASCADE ,

    CONSTRAINT fk_course_modules_module_id
        FOREIGN KEY (module_id)
            REFERENCES modules (id) ON DELETE CASCADE ,

    CONSTRAINT uq_course_modules_course_id_index UNIQUE (course_id, index),

    PRIMARY KEY (course_id, module_id)
);


CREATE TABLE module_slides
(
    module_id BIGINT,
    slide_id  BIGINT,
    index     int,

    CONSTRAINT fk_module_slides_module_id
        FOREIGN KEY (module_id)
            REFERENCES modules (id) ON DELETE CASCADE ,

    CONSTRAINT fk_module_slides_slide_id
        FOREIGN KEY (slide_id)
            REFERENCES slides (id) ON DELETE CASCADE ,


    CONSTRAINT uq_module_slides_module_id_index UNIQUE (module_id, index),

    PRIMARY KEY (module_id, slide_id)
);


CREATE TABLE enrollment
(
    account_id       BIGINT,
    course_id        BIGINT,
    completed        boolean   default false,
    current_slide_id BIGINT    default null,
    start_date       timestamp default now(),

    CONSTRAINT fk_enrollment_account_id
        FOREIGN KEY (account_id)
            REFERENCES accounts (id) ON DELETE CASCADE ,

    CONSTRAINT fk_enrollment_course_id
        FOREIGN KEY (course_id)
            REFERENCES courses (id) ON DELETE CASCADE ,

    CONSTRAINT fk_enrollment_course_current_slide_id_slide
        FOREIGN KEY (current_slide_id)
            REFERENCES slides (id) ON DELETE CASCADE ,

    PRIMARY KEY (account_id, course_id)
);

INSERT INTO courses (id, title, description, owner, organization_id)
VALUES (1,
        'ЧПУ Изучение',
        'Полное погружение в JVM, Garbage Collector и боль.',
        1,
        1);

INSERT INTO modules (id, title)
VALUES (1, 'Введение и История'),
       (2, 'Синтаксис и Типы данных');

INSERT INTO course_modules (course_id, module_id, index)
VALUES (1, 1, 1),
       (1, 2, 2);

INSERT INTO slides (id, title, description, slide_type, payload)
VALUES (1, 'ЧПУ', 'Краткий экскурс', 'TEXT',
        '{
          "content": "# ЧПУ\nТекст рассказывающий об этом...."
        }'::jsonb);

INSERT INTO slides (id, title, description, slide_type, payload)
VALUES (2, 'Лекция от Евгения Анатольевича Чепурина', 'Для чего ЧПУ', 'VIDEO_URL',
        '{
          "videoUrl": "https://youtube.com/watch?v=dQw4w9WgXcQ",
          "durationSeconds": 1200,
          "platform": "YOUTUBE"
        }'::jsonb);

INSERT INTO slides (id, title, description, slide_type, payload)
VALUES (3, 'Проверка знаний', 'Тест по первой главе', 'TEST',
        '{
          "question": "В каком году было придумано ЧПУ",
          "isMultiSelect": false,
          "options": [
            {
              "id": 1,
              "text": "1990",
              "isCorrect": false
            },
            {
              "id": 2,
              "text": "1995",
              "isCorrect": true
            },
            {
              "id": 3,
              "text": "2024",
              "isCorrect": false
            }
          ]
        }'::jsonb);

INSERT INTO module_slides (module_id, slide_id, index)
VALUES (1, 1, 1),
       (1, 2, 2);

INSERT INTO module_slides (module_id, slide_id, index)
VALUES (2, 3, 1);


SELECT setval('accounts_id_seq', (SELECT MAX(id) FROM accounts));
SELECT setval('organizations_id_seq', (SELECT MAX(id) FROM organizations));
SELECT setval('courses_id_seq', (SELECT MAX(id) FROM courses));
SELECT setval('modules_id_seq', (SELECT MAX(id) FROM modules));
SELECT setval('slides_id_seq', (SELECT MAX(id) FROM slides));

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS slides;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS courses;
