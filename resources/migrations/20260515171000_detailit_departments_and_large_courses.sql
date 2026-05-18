-- +goose Up
-- +goose StatementBegin
INSERT INTO accounts (email, password_hash, role, username)
VALUES
    ('operator@detailit.ru', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'USER'::roles, 'Оператор электрооборудования'),
    ('assembler@detailit.ru', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'USER'::roles, 'Сборщик жгутов и разъемов'),
    ('quality@detailit.ru', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'USER'::roles, 'Контролер ОТК'),
    ('engineer@detailit.ru', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'ADMIN'::roles, 'Инженер-технолог'),
    ('storekeeper@detailit.ru', '$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6', 'USER'::roles, 'Кладовщик компонентов')
ON CONFLICT (email) DO UPDATE
SET username = EXCLUDED.username,
    role = EXCLUDED.role;

UPDATE organizations
SET title = 'ООО "ДетаЛит"',
    description = 'Головной учебный контур ООО "ДетаЛит": производство электрического и электронного оборудования для автотранспортных средств, ОКВЭД 29.31.'
WHERE tag = 'detailit';

WITH departments (title, tag, description, owner_email) AS (
    VALUES
        ('ООО "ДетаЛит" — производственный участок электрооборудования', 'det-prod', 'Подразделение выпускает электрическое и электронное оборудование для автотранспортных средств: узлы, модули, сборки и производственные партии.', 'director@detailit.ru'),
        ('ООО "ДетаЛит" — участок сборки жгутов и разъемов', 'det-harness', 'Рабочая зона подготовки проводов, обжима контактов, пиновки разъемов, маркировки и финального теста цепей.', 'master@detailit.ru'),
        ('ООО "ДетаЛит" — лаборатория испытаний и диагностики', 'det-lab', 'Подразделение функциональных испытаний, электрических измерений, диагностики отказов и подтверждения параметров изделий.', 'engineer@detailit.ru'),
        ('ООО "ДетаЛит" — ОТК и метрология', 'det-qc', 'Подразделение входного, операционного и приемочного контроля, работы с несоответствиями и измерительным оборудованием.', 'master@detailit.ru'),
        ('ООО "ДетаЛит" — склад компонентов и входной контроль', 'det-store', 'Подразделение приемки электронных компонентов, ESD-хранения, маркировки, FIFO и передачи материалов в производство.', 'storekeeper@detailit.ru'),
        ('ООО "ДетаЛит" — инженерно-технологический отдел', 'det-eng', 'Подразделение технологических карт, маршрутов, анализа причин дефектов и улучшения производственного процесса.', 'engineer@detailit.ru'),
        ('ООО "ДетаЛит" — охрана труда и производственная безопасность', 'det-safety', 'Подразделение инструктажей, безопасной организации рабочих мест, допуска к операциям и фиксации происшествий.', 'director@detailit.ru')
)
INSERT INTO organizations (title, tag, description, owner_id)
SELECT d.title, d.tag, d.description, a.id
FROM departments d
JOIN accounts a ON a.email = d.owner_email
ON CONFLICT (tag) DO UPDATE
SET title = EXCLUDED.title,
    description = EXCLUDED.description,
    owner_id = EXCLUDED.owner_id;

WITH membership (tag, email) AS (
    VALUES
        ('detailit', 'director@detailit.ru'), ('detailit', 'master@detailit.ru'), ('detailit', 'worker@detailit.ru'), ('detailit', 'engineer@detailit.ru'),
        ('det-prod', 'director@detailit.ru'), ('det-prod', 'master@detailit.ru'), ('det-prod', 'operator@detailit.ru'), ('det-prod', 'worker@detailit.ru'),
        ('det-harness', 'master@detailit.ru'), ('det-harness', 'assembler@detailit.ru'), ('det-harness', 'quality@detailit.ru'), ('det-harness', 'operator@detailit.ru'),
        ('det-lab', 'engineer@detailit.ru'), ('det-lab', 'quality@detailit.ru'), ('det-lab', 'master@detailit.ru'),
        ('det-qc', 'master@detailit.ru'), ('det-qc', 'quality@detailit.ru'), ('det-qc', 'engineer@detailit.ru'),
        ('det-store', 'storekeeper@detailit.ru'), ('det-store', 'quality@detailit.ru'), ('det-store', 'master@detailit.ru'),
        ('det-eng', 'engineer@detailit.ru'), ('det-eng', 'master@detailit.ru'), ('det-eng', 'director@detailit.ru'),
        ('det-safety', 'director@detailit.ru'), ('det-safety', 'master@detailit.ru'), ('det-safety', 'operator@detailit.ru'), ('det-safety', 'assembler@detailit.ru'), ('det-safety', 'quality@detailit.ru')
)
INSERT INTO organization_accounts (organization_id, account_id)
SELECT o.id, a.id
FROM membership m
JOIN organizations o ON o.tag = m.tag
JOIN accounts a ON a.email = m.email
ON CONFLICT DO NOTHING;

CREATE OR REPLACE FUNCTION pg_temp.detailit_course(p_org_tag text, p_title text, p_description text, p_owner_email text)
RETURNS bigint AS $$
DECLARE
    v_course_id bigint;
    v_org_id bigint;
    v_owner_id bigint;
BEGIN
    SELECT id INTO v_org_id FROM organizations WHERE tag = p_org_tag;
    SELECT id INTO v_owner_id FROM accounts WHERE email = p_owner_email;

    SELECT id INTO v_course_id
    FROM courses
    WHERE title = p_title AND organization_id = v_org_id
    LIMIT 1;

    IF v_course_id IS NULL THEN
        INSERT INTO courses (title, description, owner, organization_id)
        VALUES (p_title, p_description, v_owner_id, v_org_id)
        RETURNING id INTO v_course_id;
    ELSE
        UPDATE courses
        SET description = p_description,
            owner = v_owner_id,
            organization_id = v_org_id
        WHERE id = v_course_id;
    END IF;

    RETURN v_course_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION pg_temp.detailit_module(p_course_id bigint, p_index int, p_title text)
RETURNS bigint AS $$
DECLARE
    v_module_id bigint;
BEGIN
    SELECT m.id INTO v_module_id
    FROM modules m
    JOIN course_modules cm ON cm.module_id = m.id
    WHERE cm.course_id = p_course_id AND m.title = p_title
    LIMIT 1;

    IF v_module_id IS NULL THEN
        INSERT INTO modules (title) VALUES (p_title) RETURNING id INTO v_module_id;
        INSERT INTO course_modules (course_id, module_id, index)
        VALUES (p_course_id, v_module_id, p_index)
        ON CONFLICT DO NOTHING;
    ELSE
        UPDATE course_modules
        SET index = p_index
        WHERE course_id = p_course_id AND module_id = v_module_id;
    END IF;

    RETURN v_module_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION pg_temp.detailit_slide(p_module_id bigint, p_index int, p_title text, p_description text, p_type slide_variation, p_payload jsonb)
RETURNS bigint AS $$
DECLARE
    v_slide_id bigint;
BEGIN
    SELECT s.id INTO v_slide_id
    FROM slides s
    JOIN module_slides ms ON ms.slide_id = s.id
    WHERE ms.module_id = p_module_id AND s.title = p_title
    LIMIT 1;

    IF v_slide_id IS NULL THEN
        INSERT INTO slides (title, description, slide_type, payload)
        VALUES (p_title, p_description, p_type, p_payload)
        RETURNING id INTO v_slide_id;
        INSERT INTO module_slides (module_id, slide_id, index)
        VALUES (p_module_id, v_slide_id, p_index)
        ON CONFLICT DO NOTHING;
    ELSE
        UPDATE slides
        SET description = p_description,
            slide_type = p_type,
            payload = p_payload
        WHERE id = v_slide_id;
        UPDATE module_slides
        SET index = p_index
        WHERE module_id = p_module_id AND slide_id = v_slide_id;
    END IF;

    RETURN v_slide_id;
END;
$$ LANGUAGE plpgsql;

DO $detailit$
DECLARE
    c_id bigint;
    m_id bigint;
BEGIN
    c_id := pg_temp.detailit_course(
        'det-prod',
        'Автоэлектроника и производственный маршрут изделия',
        'Большой курс для рабочих производственного участка ООО "ДетаЛит": как устроено электрическое и электронное оборудование автомобиля, почему важны жгуты, ЭБУ, датчики, прослеживаемость и стабильный маршрут партии.',
        'engineer@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'ОКВЭД 29.31 и логика изделия');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Чем занимается ДетаЛит', 'Контекст производства электрического и электронного оборудования для автотранспортных средств', 'TEXT', $json$
        {"content":"# ООО \"ДетаЛит\" и профиль 29.31\n\nКомпания обучает сотрудников работе с электрическим и электронным оборудованием для автотранспортных средств. Для рабочего это означает, что каждая операция влияет на безопасность, надежность и доказуемое качество изделия.\n\n**Что важно понимать:**\n\n- изделие состоит из электрических цепей, разъемов, электронных модулей и маркировки;\n- ошибка в маленькой операции может проявиться уже в автомобиле;\n- партия должна иметь маршрут, запись контроля и ответственного;\n- спорное изделие нельзя выпускать без решения мастера и ОТК."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Схема автоэлектроники', 'Изображение с базовой архитектурой изделия', 'FILE', $json$
        {"filename":"electronics-architecture.png","object_key":"course_files/detailit/electronics-architecture.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'От сигнала к отказу', 'Почему дефект соединения превращается в проблему у заказчика', 'TEXT', $json$
        {"content":"## Логика цепочки\n\n1. Датчик или цепь передает сигнал.\n2. Электронный блок обрабатывает состояние.\n3. Исполнительная цепь выполняет команду.\n4. Разъем, провод, пайка и контакт должны сохранить стабильность на всем пути.\n\nЕсли контакт обжат плохо, провод перепутан, плата повреждена статикой или результат испытаний не записан, качество изделия нельзя считать подтвержденным."}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Маршрут партии и контрольные точки');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Презентация: учебная сетка 29.31', 'Материал для очного разбора с мастером или технологом', 'FILE', $json$
        {"filename":"detailit-auto-electronics-training.pptx","object_key":"course_files/detailit/detailit-auto-electronics-training.pptx","mime_type":"application/vnd.openxmlformats-officedocument.presentationml.presentation"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Карта маршрута партии', 'Как партия движется между складом, производством, лабораторией и ОТК', 'TEXT', $json$
        {"content":"## Производственный маршрут\n\n- Склад принимает компоненты и проверяет маркировку.\n- Производственный участок получает задание и технологическую карту.\n- Сборщик выполняет операцию и фиксирует результат.\n- Лаборатория проверяет электрические параметры и функциональность.\n- ОТК принимает решение по партии.\n\nМаршрут нужен не для формальности: он позволяет восстановить историю изделия при внутреннем дефекте или рекламации заказчика."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: маршрут изделия', 'Контрольный вопрос по прослеживаемости', 'TEST', $json$
        {"question":"Что нужно сделать, если номер партии на изделии не совпадает с маршрутной картой?","isMultiSelect":false,"output":2,"options":[{"id":1,"text":"Передать изделие дальше, если внешний вид нормальный","isCorrect":false},{"id":2,"text":"Остановить передачу, отделить изделие и сообщить мастеру","isCorrect":true},{"id":3,"text":"Исправить номер вручную без записи","isCorrect":false}]}
        $json$::jsonb);

    c_id := pg_temp.detailit_course(
        'det-harness',
        'Жгуты, разъемы, маркировка и финальный тест цепей',
        'Практический курс для сборщиков и операторов: подготовка провода, обжим контактов, пиновка разъемов, маркировка, проверка цепей и типовые ошибки при выпуске автомобильных жгутов.',
        'master@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'Подготовка и сборка жгута');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Зачем жгуту дисциплина', 'Роль проводов, контактов и разъемов в надежности изделия', 'TEXT', $json$
        {"content":"# Жгут как часть системы\n\nЖгут связывает датчики, электронные блоки и исполнительные цепи. Ошибка в цвете провода, длине, контакте или пиновке может привести к отказу оборудования даже при исправной электронике.\n\nПеред операцией сотрудник сверяет задание, сечение, длину, цвет, маркировку, контакт и инструмент."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Схема сборки жгутов', 'Изображение процесса подготовки, обжима, пиновки и проверки', 'FILE', $json$
        {"filename":"harness-assembly.png","object_key":"course_files/detailit/harness-assembly.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Типовые ошибки сборщика', 'Что чаще всего приводит к браку', 'TEXT', $json$
        {"content":"## Ошибки, которые нельзя пропускать\n\n- перепутан контактный номер в разъеме;\n- повреждена изоляция при зачистке;\n- обжим выполнен не тем профилем инструмента;\n- маркировка нечитаема или не совпадает с заданием;\n- провод не прошел тест на непрерывность цепи.\n\nЛюбая такая ошибка блокирует передачу изделия дальше."}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Контроль и выпуск');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Контроль обжима и пиновки', 'Порядок самопроверки после операции', 'TEXT', $json$
        {"content":"## Самопроверка\n\n1. Сверить номер провода и контакт по карте.\n2. Проверить геометрию обжима и отсутствие повреждения жилы.\n3. Убедиться, что контакт зафиксирован в корпусе разъема.\n4. Проверить маркировку и направление укладки.\n5. Передать изделие на тест цепей и записать результат."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Word: регламент обучения подразделений', 'Документ с матрицей ролей, курсов и допуска', 'FILE', $json$
        {"filename":"detailit-department-training-regulation.docx","object_key":"course_files/detailit/detailit-department-training-regulation.docx","mime_type":"application/vnd.openxmlformats-officedocument.wordprocessingml.document"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: жгуты и разъемы', 'Контрольный вопрос по выпуску изделия', 'TEST', $json$
        {"question":"Что делать, если тест цепи показал обрыв после сборки разъема?","isMultiSelect":false,"output":3,"options":[{"id":1,"text":"Поставить отметку годно, если внешний вид хороший","isCorrect":false},{"id":2,"text":"Заменить маркировку и передать дальше","isCorrect":false},{"id":3,"text":"Изолировать изделие, проверить пиновку и сообщить мастеру","isCorrect":true}]}
        $json$::jsonb);

    c_id := pg_temp.detailit_course(
        'det-eng',
        'ESD, пайка и монтаж электронных компонентов',
        'Расширенный курс для технологов, операторов и мастеров: организация ESD-поста, защита компонентов, качество пайки, работа с платами, визуальный контроль и предотвращение скрытых отказов.',
        'engineer@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'ESD-пост и защита компонентов');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Почему ESD опасен', 'Скрытые повреждения электронных компонентов', 'TEXT', $json$
        {"content":"# ESD и скрытый отказ\n\nЭлектростатический разряд может повредить компонент так, что изделие пройдет первичный запуск, но откажет позже. Поэтому ESD-дисциплина обязательна до касания платы или электронного модуля.\n\nСотрудник проверяет браслет, коврик, заземление, тару и состояние рабочего места перед сменой."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Изображение: ESD-пост', 'Наглядная схема рабочего места монтажа', 'FILE', $json$
        {"filename":"esd-workstation.png","object_key":"course_files/detailit/esd-workstation.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'PDF: стандарт ESD и пайки', 'Полный документ с нормами поста, пайки и прослеживаемости', 'FILE', $json$
        {"filename":"detailit-esd-soldering-standard.pdf","object_key":"course_files/detailit/detailit-esd-soldering-standard.pdf","mime_type":"application/pdf"}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Пайка, монтаж и визуальный контроль');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Признаки качественной пайки', 'Что проверять перед передачей платы дальше', 'TEXT', $json$
        {"content":"## Визуальный контроль пайки\n\nКачественное соединение соответствует техкарте, не имеет мостиков, шариков припоя, трещин, перегрева площадки, непропая и загрязнений. Если признак спорный, сотрудник не принимает решение один: плата отделяется и передается мастеру или ОТК."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Алгоритм спорного монтажа', 'Порядок действий при сомнительной операции', 'TEXT', $json$
        {"content":"1. Остановить передачу платы дальше.\n2. Сохранить изделие в ESD-таре.\n3. Зафиксировать номер партии и операцию.\n4. Сообщить мастеру или технологу.\n5. Выполнить доработку только по разрешенной процедуре.\n6. Провести повторный контроль и записать результат."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: ESD и пайка', 'Контрольный вопрос по рабочему месту', 'TEST', $json$
        {"question":"Можно ли начинать монтаж платы, если ESD-браслет не прошел проверку?","isMultiSelect":false,"output":1,"options":[{"id":1,"text":"Нет, монтаж нужно остановить до исправления защиты","isCorrect":true},{"id":2,"text":"Да, если операция короткая","isCorrect":false},{"id":3,"text":"Да, если плата лежит на столе отдельно","isCorrect":false}]}
        $json$::jsonb);

    c_id := pg_temp.detailit_course(
        'det-lab',
        'Испытания, ОТК и трассируемость партии',
        'Курс для лаборатории, ОТК, мастеров и инженеров: функциональные испытания, электрические измерения, входной контроль компонентов, оформление несоответствий и доказательство качества партии.',
        'master@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'Испытания и измерения');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Зачем нужны испытания', 'Контроль качества электронного изделия перед выпуском', 'TEXT', $json$
        {"content":"# Испытания подтверждают качество\n\nДля продукции 29.31 недостаточно внешнего осмотра. Нужно подтвердить электрические параметры, функциональность, маркировку, связь с партией и соответствие методике.\n\nЕсли результат не записан, качество нельзя доказать даже при фактически исправном изделии."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Схема трассируемости', 'Изображение связей между партией, измерением и решением ОТК', 'FILE', $json$
        {"filename":"testing-traceability.png","object_key":"course_files/detailit/testing-traceability.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Методика измерения', 'Минимальные правила записи результата', 'TEXT', $json$
        {"content":"## Что должно быть в записи\n\n- номер партии и изделия;\n- проверяемый параметр и допустимый диапазон;\n- прибор или стенд;\n- результат измерения;\n- исполнитель, дата, смена;\n- решение: годно, доработка, изоляция, брак.\n\nЗапись без метода или допуска не помогает принять решение."}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Несоответствия и рекламации');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Как работать с несоответствием', 'Порядок изоляции спорного изделия', 'TEXT', $json$
        {"content":"1. Остановить передачу изделия дальше.\n2. Отделить спорную единицу или партию.\n3. Зафиксировать признак, операцию, время и исполнителя.\n4. Сообщить мастеру, ОТК или инженеру.\n5. Не смешивать спорные изделия с годными.\n6. Выполнить решение: доработка, повторный контроль, списание или выпуск по разрешению."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Презентация для разбора дефектов', 'Материал для обучения лаборатории и ОТК', 'FILE', $json$
        {"filename":"detailit-auto-electronics-training.pptx","object_key":"course_files/detailit/detailit-auto-electronics-training.pptx","mime_type":"application/vnd.openxmlformats-officedocument.presentationml.presentation"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: решение ОТК', 'Контрольный вопрос по несоответствиям', 'TEST', $json$
        {"question":"Что делать с изделием, которое не прошло функциональный тест?","isMultiSelect":false,"output":2,"options":[{"id":1,"text":"Вернуть в годную тару до конца смены","isCorrect":false},{"id":2,"text":"Изолировать, записать результат и передать на решение ОТК/мастера","isCorrect":true},{"id":3,"text":"Стереть результат и повторить тест без записи","isCorrect":false}]}
        $json$::jsonb);
END;
$detailit$;

SELECT setval('accounts_id_seq', (SELECT MAX(id) FROM accounts));
SELECT setval('organizations_id_seq', (SELECT MAX(id) FROM organizations));
SELECT setval('courses_id_seq', (SELECT MAX(id) FROM courses));
SELECT setval('modules_id_seq', (SELECT MAX(id) FROM modules));
SELECT setval('slides_id_seq', (SELECT MAX(id) FROM slides));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM courses
WHERE title IN (
    'Автоэлектроника и производственный маршрут изделия',
    'Жгуты, разъемы, маркировка и финальный тест цепей',
    'ESD, пайка и монтаж электронных компонентов',
    'Испытания, ОТК и трассируемость партии'
);

DELETE FROM slides
WHERE title IN (
    'Чем занимается ДетаЛит',
    'Схема автоэлектроники',
    'От сигнала к отказу',
    'Презентация: учебная сетка 29.31',
    'Карта маршрута партии',
    'Проверка: маршрут изделия',
    'Зачем жгуту дисциплина',
    'Схема сборки жгутов',
    'Типовые ошибки сборщика',
    'Контроль обжима и пиновки',
    'Word: регламент обучения подразделений',
    'Проверка: жгуты и разъемы',
    'Почему ESD опасен',
    'Изображение: ESD-пост',
    'PDF: стандарт ESD и пайки',
    'Признаки качественной пайки',
    'Алгоритм спорного монтажа',
    'Проверка: ESD и пайка',
    'Зачем нужны испытания',
    'Схема трассируемости',
    'Методика измерения',
    'Как работать с несоответствием',
    'Презентация для разбора дефектов',
    'Проверка: решение ОТК'
);

DELETE FROM modules
WHERE title IN (
    'ОКВЭД 29.31 и логика изделия',
    'Маршрут партии и контрольные точки',
    'Подготовка и сборка жгута',
    'Контроль и выпуск',
    'ESD-пост и защита компонентов',
    'Пайка, монтаж и визуальный контроль',
    'Испытания и измерения',
    'Несоответствия и рекламации'
);

DELETE FROM organization_accounts
WHERE organization_id IN (SELECT id FROM organizations WHERE tag IN ('det-prod', 'det-harness', 'det-lab', 'det-qc', 'det-store', 'det-eng', 'det-safety'));

DELETE FROM organizations
WHERE tag IN ('det-prod', 'det-harness', 'det-lab', 'det-qc', 'det-store', 'det-eng', 'det-safety');

DELETE FROM accounts
WHERE email IN ('operator@detailit.ru', 'assembler@detailit.ru', 'quality@detailit.ru', 'engineer@detailit.ru', 'storekeeper@detailit.ru');
-- +goose StatementEnd
