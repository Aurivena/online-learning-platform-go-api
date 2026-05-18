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
            REFERENCES accounts (id) ON DELETE CASCADE,

    CONSTRAINT fk_course_organization_id
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
            REFERENCES courses (id) ON DELETE CASCADE,

    CONSTRAINT fk_course_modules_module_id
        FOREIGN KEY (module_id)
            REFERENCES modules (id) ON DELETE CASCADE,

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
            REFERENCES modules (id) ON DELETE CASCADE,

    CONSTRAINT fk_module_slides_slide_id
        FOREIGN KEY (slide_id)
            REFERENCES slides (id) ON DELETE CASCADE,

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
            REFERENCES accounts (id) ON DELETE CASCADE,

    CONSTRAINT fk_enrollment_course_id
        FOREIGN KEY (course_id)
            REFERENCES courses (id) ON DELETE CASCADE,

    CONSTRAINT fk_enrollment_current_slide_id
        FOREIGN KEY (current_slide_id)
            REFERENCES slides (id) ON DELETE SET NULL,

    PRIMARY KEY (account_id, course_id)
);

CREATE INDEX idx_courses_organization_id ON courses (organization_id);
CREATE INDEX idx_course_modules_course_id ON course_modules (course_id);
CREATE INDEX idx_module_slides_module_id ON module_slides (module_id);

INSERT INTO courses (id, title, description, owner, organization_id)
VALUES (1,
        'Безопасность и вводный инструктаж на производстве',
        'Обязательный курс для сотрудников ООО "ДетаЛит": рабочая зона, средства защиты, действия при отклонениях и правила фиксации происшествий.',
        1,
        1),
       (2,
        'Технологический маршрут: от формы до готовой отливки',
        'Практический курс по этапам изготовления: подготовка формы, плавка, заливка, охлаждение, передача партии в ОТК.',
        2,
        1),
       (3,
        'Контроль качества и работа с несоответствиями',
        'Курс для мастеров, операторов и сотрудников ОТК: дефекты, измерения, документы партии и порядок изоляции брака.',
        2,
        1);

INSERT INTO modules (id, title)
VALUES (1, 'Вводный инструктаж'),
       (2, 'Отклонения и аварийные ситуации'),
       (3, 'Маршрут изготовления'),
       (4, 'Плавка, заливка и охлаждение'),
       (5, 'Контроль отливок'),
       (6, 'Приемка партии и фиксация результата');

INSERT INTO course_modules (course_id, module_id, index)
VALUES (1, 1, 1),
       (1, 2, 2),
       (2, 3, 1),
       (2, 4, 2),
       (3, 5, 1),
       (3, 6, 2);

INSERT INTO slides (id, title, description, slide_type, payload)
VALUES (1,
        'Зачем нужен вводный инструктаж',
        'Контекст курса и ожидаемое поведение сотрудника на участке',
        'TEXT',
        $json$
        {
          "content": "# Вводный инструктаж ООО \"ДетаЛит\"\n\nКурс помогает сотруднику безопасно войти в смену и понимать, какие действия обязательны до начала работы.\n\n**После прохождения вы сможете:**\n\n- проверить рабочую зону и средства защиты;\n- распознать опасное отклонение;\n- правильно сообщить мастеру о событии;\n- отличить учебную ситуацию от реального риска для партии и оборудования.\n\nЗапомните базовое правило: если операция небезопасна или результат нельзя проверить документально, работу нужно остановить и уточнить у мастера."
        }
        $json$::jsonb),
       (2,
        'Изображение: безопасная рабочая зона',
        'Наглядная памятка по рабочей зоне, средствам защиты и сигналам отклонений',
        'FILE',
        $json$
        {
          "filename": "safety-zone.png",
          "object_key": "course_files/detailit/safety-zone.png",
          "mime_type": "image/png"
        }
        $json$::jsonb),
       (3,
        'Презентация: безопасность, процесс и качество',
        'Краткая презентация для вводного занятия и повторения ключевых правил',
        'FILE',
        $json$
        {
          "filename": "detailit-production-intro.pptx",
          "object_key": "course_files/detailit/detailit-production-intro.pptx",
          "mime_type": "application/vnd.openxmlformats-officedocument.presentationml.presentation"
        }
        $json$::jsonb),
       (4,
        'Алгоритм действий при отклонении',
        'Что делать, если на участке замечена опасность или дефект',
        'TEXT',
        $json$
        {
          "content": "## Алгоритм сотрудника\n\n1. Остановите операцию, если продолжение может привести к травме, браку или повреждению оборудования.\n2. Обозначьте опасную зону и предупредите коллег рядом.\n3. Сообщите мастеру смены: участок, время, операция, признак отклонения.\n4. Отделите спорную деталь или партию, если отклонение связано с качеством.\n5. Запишите событие в журнал смены или передайте данные ответственному.\n\nНе пытайтесь самостоятельно обходить блокировки, отключать защиту или выпускать спорное изделие дальше по маршруту."
        }
        $json$::jsonb),
       (5,
        'Word-документ: карта ежедневного инструктажа',
        'Документ с чек-листами смены, таблицей контрольных точек и листом самопроверки',
        'FILE',
        $json$
        {
          "filename": "detailit-instruction-card.docx",
          "object_key": "course_files/detailit/detailit-instruction-card.docx",
          "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        }
        $json$::jsonb),
       (6,
        'Проверка: безопасный старт смены',
        'Контрольный вопрос по вводному инструктажу',
        'TEST',
        $json$
        {
          "question": "Что нужно сделать, если сотрудник заметил перегрев оборудования перед запуском операции?",
          "isMultiSelect": false,
          "output": 2,
          "options": [
            { "id": 1, "text": "Запустить операцию и проверить позже", "isCorrect": false },
            { "id": 2, "text": "Остановить старт, предупредить коллег и сообщить мастеру", "isCorrect": true },
            { "id": 3, "text": "Отключить защиту, чтобы не мешала работе", "isCorrect": false }
          ]
        }
        $json$::jsonb),
       (7,
        'Полный маршрут партии',
        'Как партия проходит путь от задания смены до ОТК',
        'TEXT',
        $json$
        {
          "content": "# Маршрут изготовления отливки\n\nНа ООО \"ДетаЛит\" каждая партия должна проходить одинаковый контролируемый путь:\n\n1. Получение задания смены и маршрутной карты.\n2. Подготовка формы и проверка маркировки.\n3. Плавка с контролем температуры и состава.\n4. Заливка без нарушения защит и технологической карты.\n5. Охлаждение с выдержкой по регламенту.\n6. Очистка, первичный осмотр и передача в ОТК.\n\nЕсли этап пропущен, качество партии нельзя доказать документально."
        }
        $json$::jsonb),
       (8,
        'Изображение: маршрут изготовления',
        'Схема производственного процесса от формы до контроля качества',
        'FILE',
        $json$
        {
          "filename": "process-flow.png",
          "object_key": "course_files/detailit/process-flow.png",
          "mime_type": "image/png"
        }
        $json$::jsonb),
       (9,
        'Презентация для технологического разбора',
        'Материал можно использовать на очном инструктаже по маршруту партии',
        'FILE',
        $json$
        {
          "filename": "detailit-production-intro.pptx",
          "object_key": "course_files/detailit/detailit-production-intro.pptx",
          "mime_type": "application/vnd.openxmlformats-officedocument.presentationml.presentation"
        }
        $json$::jsonb),
       (10,
        'Критичные параметры процесса',
        'Какие параметры нельзя оставлять без контроля',
        'TEXT',
        $json$
        {
          "content": "## Критичные параметры\n\n- **Температура плавки:** фиксируется по технологической карте, отклонение согласуется с мастером.\n- **Состояние формы:** сколы, загрязнения и перекосы повышают риск брака.\n- **Скорость заливки:** резкие изменения могут привести к раковинам и непроливу.\n- **Время охлаждения:** досрочное вскрытие формы ухудшает геометрию и поверхность.\n- **Маркировка:** без связи с партией изделие нельзя уверенно принять."
        }
        $json$::jsonb),
       (11,
        'Проверка: маршрут партии',
        'Контрольный вопрос по технологическому маршруту',
        'TEST',
        $json$
        {
          "question": "Почему нельзя пропускать запись результата контроля в маршрутной карте?",
          "isMultiSelect": false,
          "output": 3,
          "options": [
            { "id": 1, "text": "Потому что так быстрее закончится смена", "isCorrect": false },
            { "id": 2, "text": "Потому что запись нужна только бухгалтерии", "isCorrect": false },
            { "id": 3, "text": "Потому что без записи нельзя подтвердить качество и прослеживаемость партии", "isCorrect": true }
          ]
        }
        $json$::jsonb),
       (12,
        'Что проверяет сотрудник ОТК',
        'Базовые признаки годного изделия и брака',
        'TEXT',
        $json$
        {
          "content": "# Контроль качества отливок\n\nКонтроль не сводится к поиску одного дефекта. Сотрудник сверяет изделие по четырем группам признаков:\n\n- поверхность: трещины, раковины, непроливы, подгар;\n- геометрия: базовые размеры, плоскостность, отверстия, посадочные зоны;\n- маркировка: партия, смена, маршрутная карта;\n- решение: годно, доработка, изоляция или брак.\n\nСпорное изделие нельзя выпускать дальше без решения мастера и ОТК."
        }
        $json$::jsonb),
       (13,
        'Изображение: контрольные точки качества',
        'Наглядная схема проверки поверхности, геометрии, маркировки и решения по изделию',
        'FILE',
        $json$
        {
          "filename": "quality-checkpoints.png",
          "object_key": "course_files/detailit/quality-checkpoints.png",
          "mime_type": "image/png"
        }
        $json$::jsonb),
       (14,
        'PDF: стандарт визуального и измерительного контроля',
        'Регламент контроля партии после охлаждения, очистки и первичной обработки',
        'FILE',
        $json$
        {
          "filename": "detailit-quality-standard.pdf",
          "object_key": "course_files/detailit/detailit-quality-standard.pdf",
          "mime_type": "application/pdf"
        }
        $json$::jsonb),
       (15,
        'Как оформлять несоответствие',
        'Порядок изоляции партии и передачи информации мастеру',
        'TEXT',
        $json$
        {
          "content": "## Несоответствие: минимальная запись\n\nВ журнале или карточке события должны быть:\n\n1. номер партии и наименование детали;\n2. участок, операция, смена и время обнаружения;\n3. вид отклонения: поверхность, размер, маркировка, документ;\n4. действие: изоляция, повторное измерение, доработка, решение ОТК;\n5. фамилия сотрудника, передавшего информацию мастеру.\n\nЦель записи — не наказание, а прослеживаемость и снижение повторяемости дефекта."
        }
        $json$::jsonb),
       (16,
        'Word-документ для самопроверки',
        'Карта инструктажа подходит как приложение к курсу по качеству',
        'FILE',
        $json$
        {
          "filename": "detailit-instruction-card.docx",
          "object_key": "course_files/detailit/detailit-instruction-card.docx",
          "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        }
        $json$::jsonb),
       (17,
        'Проверка: решение по спорной детали',
        'Контрольный вопрос по курсу качества',
        'TEST',
        $json$
        {
          "question": "Что делать со спорной деталью, если дефект неочевиден?",
          "isMultiSelect": false,
          "output": 1,
          "options": [
            { "id": 1, "text": "Отделить деталь или партию и передать решение мастеру и ОТК", "isCorrect": true },
            { "id": 2, "text": "Отправить дальше, если визуально выглядит нормально", "isCorrect": false },
            { "id": 3, "text": "Стереть маркировку, чтобы не задерживать маршрут", "isCorrect": false }
          ]
        }
        $json$::jsonb);

INSERT INTO module_slides (module_id, slide_id, index)
VALUES (1, 1, 1),
       (1, 2, 2),
       (1, 3, 3),
       (2, 4, 1),
       (2, 5, 2),
       (2, 6, 3),
       (3, 7, 1),
       (3, 8, 2),
       (3, 9, 3),
       (4, 10, 1),
       (4, 11, 2),
       (5, 12, 1),
       (5, 13, 2),
       (5, 14, 3),
       (6, 15, 1),
       (6, 16, 2),
       (6, 17, 3);

INSERT INTO enrollment (account_id, course_id, completed, current_slide_id, start_date)
VALUES (3, 1, false, 1, now()),
       (3, 2, false, 7, now()),
       (3, 3, false, 12, now()),
       (2, 1, false, 1, now()),
       (2, 2, false, 7, now()),
       (2, 3, false, 12, now());

SELECT setval('accounts_id_seq', (SELECT MAX(id) FROM accounts));
SELECT setval('organizations_id_seq', (SELECT MAX(id) FROM organizations));
SELECT setval('courses_id_seq', (SELECT MAX(id) FROM courses));
SELECT setval('modules_id_seq', (SELECT MAX(id) FROM modules));
SELECT setval('slides_id_seq', (SELECT MAX(id) FROM slides));

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS enrollment;
DROP TABLE IF EXISTS module_slides;
DROP TABLE IF EXISTS course_modules;
DROP TABLE IF EXISTS slides;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS courses;
DROP TYPE IF EXISTS slide_variation;
