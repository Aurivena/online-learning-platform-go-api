-- +goose Up
-- +goose StatementBegin
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
            owner = v_owner_id
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
        'det-store',
        'Входной контроль компонентов, ESD-хранение и FIFO',
        'Курс для склада компонентов и входного контроля ООО "ДетаЛит": приемка электронных компонентов, проверка документов, ESD-упаковка, маркировка, FIFO и выдача материалов в производство.',
        'storekeeper@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'Приемка компонентов');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Склад как первая контрольная точка', 'Почему входной контроль влияет на качество изделия', 'TEXT', $json$
        {"content":"# Входной контроль\n\nСклад не просто хранит материалы. Он защищает производство от неправильных компонентов, поврежденной упаковки, смешения партий и ESD-рисков.\n\nПри приемке проверяются поставщик, документ, маркировка, количество, состояние упаковки, срок хранения и соответствие заказу."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'ESD-упаковка и хранение', 'Связь склада с защитой электронных компонентов', 'FILE', $json$
        {"filename":"esd-workstation.png","object_key":"course_files/detailit/esd-workstation.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Ошибки входного контроля', 'Что нельзя пропускать при приемке', 'TEXT', $json$
        {"content":"## Красные флаги\n\n- поврежденная или обычная не-ESD упаковка;\n- нечитаемая маркировка партии;\n- несоответствие количества документам;\n- смешение разных ревизий компонента;\n- нет связи с заказом или сертификатом.\n\nТакой материал не передается в производство без решения ответственного."}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Выдача в производство');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'FIFO и прослеживаемость', 'Как склад поддерживает маршрут партии', 'TEXT', $json$
        {"content":"Материалы выдаются по FIFO или по отдельному распоряжению технолога. В записи должны быть номер партии компонента, количество, получатель, подразделение, дата и основание выдачи. Это позволяет связать готовое изделие с конкретной поставкой."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Документ: регламент обучения', 'Файл с матрицей ролей и допуска', 'FILE', $json$
        {"filename":"detailit-department-training-regulation.docx","object_key":"course_files/detailit/detailit-department-training-regulation.docx","mime_type":"application/vnd.openxmlformats-officedocument.wordprocessingml.document"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: склад компонентов', 'Контрольный вопрос по входному контролю', 'TEST', $json$
        {"question":"Что делать с компонентами в поврежденной ESD-упаковке?","isMultiSelect":false,"output":2,"options":[{"id":1,"text":"Выдать быстрее, пока упаковка не порвалась сильнее","isCorrect":false},{"id":2,"text":"Изолировать, зафиксировать отклонение и дождаться решения ответственного","isCorrect":true},{"id":3,"text":"Пересыпать в обычный пакет и промаркировать вручную","isCorrect":false}]}
        $json$::jsonb);

    c_id := pg_temp.detailit_course(
        'det-qc',
        'ОТК, метрология и работа с рекламациями',
        'Курс для контролеров, мастеров и инженеров: операционный контроль, измерительные приборы, критерии приемки, оформление несоответствий и анализ рекламаций по электронному оборудованию.',
        'master@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'Операционный контроль и метрология');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Роль ОТК в ДетаЛит', 'Что контролер подтверждает перед выпуском партии', 'TEXT', $json$
        {"content":"# ОТК подтверждает качество\n\nКонтролер проверяет не только изделие, но и доказательства: маршрут, измерение, прибор, допуск, подпись, дату и решение по партии. Без этих данных выпуск становится спорным."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Схема испытаний и трассируемости', 'Визуальная карта контроля партии', 'FILE', $json$
        {"filename":"testing-traceability.png","object_key":"course_files/detailit/testing-traceability.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Прибор и методика', 'Минимальный стандарт измерения', 'TEXT', $json$
        {"content":"Измерение считается полезным, если указаны прибор, дата поверки или допуска, методика, параметр, диапазон, фактический результат и решение. Если прибор не подходит или методика не определена, результат нельзя использовать для приемки."}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Несоответствия и рекламации');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'От внутреннего дефекта к корректирующему действию', 'Как ОТК запускает улучшение процесса', 'TEXT', $json$
        {"content":"1. Зафиксировать дефект и изолировать изделие.\n2. Определить партию, операцию и сотрудника.\n3. Передать информацию мастеру и инженеру.\n4. Провести разбор причины.\n5. Назначить действие: обучение, изменение карты, проверка инструмента, доработка или списание.\n6. Проверить, что действие реально снизило риск повторения."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'PDF: стандарт контроля', 'Документ со стандартом визуального и измерительного контроля', 'FILE', $json$
        {"filename":"detailit-quality-standard.pdf","object_key":"course_files/detailit/detailit-quality-standard.pdf","mime_type":"application/pdf"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: ОТК и метрология', 'Контрольный вопрос по измерениям', 'TEST', $json$
        {"question":"Можно ли принять изделие по измерению, если неизвестна методика и допуск?","isMultiSelect":false,"output":1,"options":[{"id":1,"text":"Нет, результат нельзя использовать для приемки","isCorrect":true},{"id":2,"text":"Да, если контролер опытный","isCorrect":false},{"id":3,"text":"Да, если изделие выглядит исправным","isCorrect":false}]}
        $json$::jsonb);

    c_id := pg_temp.detailit_course(
        'det-safety',
        'Безопасный допуск к операциям 29.31',
        'Курс для рабочих, мастеров и руководителей подразделений: безопасный старт смены, допуск к операциям с электрооборудованием, ESD-постами, испытательными стендами и действия при опасных отклонениях.',
        'director@detailit.ru'
    );
    m_id := pg_temp.detailit_module(c_id, 1, 'Допуск к смене');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Безопасность в производстве 29.31', 'Особенности рабочих мест с электроникой и испытаниями', 'TEXT', $json$
        {"content":"# Безопасный допуск\n\nПеред началом смены сотрудник проверяет рабочее место, СИЗ, исправность оборудования, отсутствие посторонних предметов, доступность аварийных средств и актуальность задания. Если операция непонятна или небезопасна, ее нельзя начинать без мастера."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Безопасная рабочая зона', 'Памятка по рабочей зоне и сигналам отклонений', 'FILE', $json$
        {"filename":"safety-zone.png","object_key":"course_files/detailit/safety-zone.png","mime_type":"image/png"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Опасные отклонения', 'Что требует немедленной остановки операции', 'TEXT', $json$
        {"content":"Операцию нужно остановить при запахе гари, перегреве, искрении, необычном шуме, повреждении кабеля, отказе блокировки, отсутствии СИЗ, нарушении ESD-защиты, подозрении на короткое замыкание или невозможности подтвердить маршрут партии."}
        $json$::jsonb);
    m_id := pg_temp.detailit_module(c_id, 2, 'Действия при событии');
    PERFORM pg_temp.detailit_slide(m_id, 1, 'Алгоритм остановки', 'Порядок действий сотрудника при риске', 'TEXT', $json$
        {"content":"1. Остановить операцию безопасным способом.\n2. Предупредить коллег рядом.\n3. Обозначить опасную зону.\n4. Сообщить мастеру: участок, операция, признак, время.\n5. Не включать оборудование до разрешения.\n6. Зафиксировать событие в журнале смены."}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 2, 'Презентация: безопасность и качество', 'Вводная презентация для инструктажа', 'FILE', $json$
        {"filename":"detailit-production-intro.pptx","object_key":"course_files/detailit/detailit-production-intro.pptx","mime_type":"application/vnd.openxmlformats-officedocument.presentationml.presentation"}
        $json$::jsonb);
    PERFORM pg_temp.detailit_slide(m_id, 3, 'Проверка: безопасный допуск', 'Контрольный вопрос по опасному отклонению', 'TEST', $json$
        {"question":"Что делать, если на испытательном стенде появился запах гари?","isMultiSelect":false,"output":3,"options":[{"id":1,"text":"Дождаться конца теста","isCorrect":false},{"id":2,"text":"Открыть корпус и искать причину одному","isCorrect":false},{"id":3,"text":"Остановить операцию, предупредить коллег и сообщить мастеру","isCorrect":true}]}
        $json$::jsonb);
END;
$detailit$;

SELECT setval('courses_id_seq', (SELECT MAX(id) FROM courses));
SELECT setval('modules_id_seq', (SELECT MAX(id) FROM modules));
SELECT setval('slides_id_seq', (SELECT MAX(id) FROM slides));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM courses
WHERE title IN (
    'Входной контроль компонентов, ESD-хранение и FIFO',
    'ОТК, метрология и работа с рекламациями',
    'Безопасный допуск к операциям 29.31'
);

DELETE FROM slides
WHERE title IN (
    'Склад как первая контрольная точка',
    'ESD-упаковка и хранение',
    'Ошибки входного контроля',
    'FIFO и прослеживаемость',
    'Документ: регламент обучения',
    'Проверка: склад компонентов',
    'Роль ОТК в ДетаЛит',
    'Схема испытаний и трассируемости',
    'Прибор и методика',
    'От внутреннего дефекта к корректирующему действию',
    'PDF: стандарт контроля',
    'Проверка: ОТК и метрология',
    'Безопасность в производстве 29.31',
    'Безопасная рабочая зона',
    'Опасные отклонения',
    'Алгоритм остановки',
    'Презентация: безопасность и качество',
    'Проверка: безопасный допуск'
);

DELETE FROM modules
WHERE title IN (
    'Приемка компонентов',
    'Выдача в производство',
    'Операционный контроль и метрология',
    'Несоответствия и рекламации',
    'Допуск к смене',
    'Действия при событии'
);
-- +goose StatementEnd
