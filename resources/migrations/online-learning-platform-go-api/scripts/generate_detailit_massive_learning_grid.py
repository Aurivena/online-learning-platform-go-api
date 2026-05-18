from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "resources" / "migrations" / "20260515183000_detailit_massive_learning_grid.sql"
PASSWORD_HASH = "$2a$12$0/3t0IVWUy79ICq7iZL/KehehrMhk1WvHpIGBebTNuRG8mKl9MxF6"

VIDEOS = [
    "https://www.youtube.com/watch?v=VD009jiZreo",
    "https://www.youtube.com/watch?v=OSUvXC1pACA",
    "https://www.youtube.com/watch?v=kCp5yYjo9zE",
    "https://www.youtube.com/watch?v=V9RLc9EX1so",
]

FILES = [
    ("electronics-architecture.png", "image/png"),
    ("harness-assembly.png", "image/png"),
    ("esd-workstation.png", "image/png"),
    ("testing-traceability.png", "image/png"),
    ("detailit-auto-electronics-training.pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation"),
    ("detailit-department-training-regulation.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"),
    ("detailit-esd-soldering-standard.pdf", "application/pdf"),
    ("detailit-quality-standard.pdf", "application/pdf"),
]

ACCOUNTS = [
    ("linelead@detailit.ru", "Старший смены линии", "ADMIN"),
    ("tester@detailit.ru", "Оператор испытательного стенда", "USER"),
    ("smt@detailit.ru", "Оператор SMT-линии", "USER"),
    ("repair@detailit.ru", "Специалист по ремонту и доработке", "USER"),
    ("metrolog@detailit.ru", "Метролог", "ADMIN"),
    ("planner@detailit.ru", "Планировщик производства", "ADMIN"),
    ("pack@detailit.ru", "Упаковщик готовой продукции", "USER"),
    ("maintenance@detailit.ru", "Слесарь-наладчик оборудования", "USER"),
]

DEPARTMENTS = [
    ("det-smt", "ООО \"ДетаЛит\" — SMT-линия и монтаж плат", "Поверхностный монтаж, подготовка плат, печать пасты, установка компонентов, оплавление и визуальный контроль.", "engineer@detailit.ru"),
    ("det-pcb", "ООО \"ДетаЛит\" — участок подготовки печатных плат", "Приемка, маркировка, очистка, хранение и подготовка печатных плат к монтажу.", "smt@detailit.ru"),
    ("det-calib", "ООО \"ДетаЛит\" — метрология и калибровка стендов", "Контроль измерительных приборов, стендов, калибровки, методик и пригодности результатов.", "metrolog@detailit.ru"),
    ("det-repair", "ООО \"ДетаЛит\" — участок ремонта и доработки", "Доработка изделий, анализ дефектов, безопасная перепайка, повторный контроль и возврат в маршрут.", "repair@detailit.ru"),
    ("det-pack", "ООО \"ДетаЛит\" — упаковка и отгрузка", "Финальная упаковка, маркировка, ESD-защита, комплектация и передача готовой продукции заказчику.", "pack@detailit.ru"),
    ("det-log", "ООО \"ДетаЛит\" — внутренняя логистика", "Перемещение партий между складом, производством, лабораторией, ОТК и упаковкой без потери прослеживаемости.", "planner@detailit.ru"),
    ("det-plan", "ООО \"ДетаЛит\" — планирование производства", "План смен, загрузка линий, маршрут партий, приоритеты заказов и управление узкими местами.", "planner@detailit.ru"),
    ("det-train", "ООО \"ДетаЛит\" — учебный центр и наставники", "Подготовка рабочих, матрица допуска, наставничество, проверки знаний и повторные инструктажи.", "director@detailit.ru"),
    ("det-maint", "ООО \"ДетаЛит\" — обслуживание оборудования", "Профилактика, диагностика, безопасная остановка, ремонт и запуск производственного оборудования.", "maintenance@detailit.ru"),
    ("det-line1", "ООО \"ДетаЛит\" — линия сборки электронных модулей 1", "Серийная сборка электронных модулей, контроль операций и выпуск партий.", "linelead@detailit.ru"),
    ("det-line2", "ООО \"ДетаЛит\" — линия сборки электронных модулей 2", "Параллельная производственная линия для стабильного выпуска и обучения сменных бригад.", "linelead@detailit.ru"),
    ("det-proto", "ООО \"ДетаЛит\" — опытный участок и прототипирование", "Пробные партии, технологические эксперименты, подтверждение маршрута и подготовка серийного запуска.", "engineer@detailit.ru"),
]

COURSES = [
    ("detailit", "Безопасность и вводный инструктаж на производстве", "director@detailit.ru", ["Вводный инструктаж", "Отклонения и аварийные ситуации", "Электробезопасность рабочего места", "Пожарные риски и эвакуация", "Финальная проверка допуска"]),
    ("detailit", "Технологический маршрут: от формы до готовой отливки", "master@detailit.ru", ["Маршрут изготовления", "Плавка, заливка и охлаждение", "Передача партии между участками", "Документы и прослеживаемость", "Контрольные вопросы маршрута"]),
    ("detailit", "Контроль качества и работа с несоответствиями", "master@detailit.ru", ["Контроль отливок", "Приемка партии и фиксация результата", "Изоляция спорной продукции", "Разбор причины дефекта", "Корректирующие действия"]),
    ("det-prod", "Автоэлектроника и производственный маршрут изделия", "engineer@detailit.ru", ["ОКВЭД 29.31 и логика изделия", "Маршрут партии и контрольные точки", "Критичные операции производства", "Рабочая документация", "Итоговая производственная практика"]),
    ("det-harness", "Жгуты, разъемы, маркировка и финальный тест цепей", "master@detailit.ru", ["Подготовка и сборка жгута", "Контроль и выпуск", "Обжим и пиновка", "Маркировка и укладка", "Финальная проверка цепей"]),
    ("det-eng", "ESD, пайка и монтаж электронных компонентов", "engineer@detailit.ru", ["ESD-пост и защита компонентов", "Пайка, монтаж и визуальный контроль", "Технологическая карта", "Доработка и повторный контроль", "Контроль знаний по монтажу"]),
    ("det-lab", "Испытания, ОТК и трассируемость партии", "master@detailit.ru", ["Испытания и измерения", "Несоответствия и рекламации", "Функциональный стенд", "Запись результатов", "Решение по партии"]),
    ("det-store", "Входной контроль компонентов, ESD-хранение и FIFO", "storekeeper@detailit.ru", ["Приемка компонентов", "Выдача в производство", "ESD-хранение", "FIFO и партии", "Проверка складского допуска"]),
    ("det-qc", "ОТК, метрология и работа с рекламациями", "master@detailit.ru", ["Операционный контроль и метрология", "Несоответствия и рекламации", "Измерительное оборудование", "Критерии приемки", "Анализ рекламаций"]),
    ("det-safety", "Безопасный допуск к операциям 29.31", "director@detailit.ru", ["Допуск к смене", "Действия при событии", "Риски электроники и стендов", "Работа с СИЗ", "Итоговый инструктаж"]),
]

COURSE_THEMES = [
    ("det-prod", "Серийный выпуск электронного узла автомобиля", "operator@detailit.ru"),
    ("det-prod", "Рабочее место оператора производственной линии", "linelead@detailit.ru"),
    ("det-safety", "Опасные энергии, блокировки и безопасная остановка", "director@detailit.ru"),
    ("det-safety", "Первая реакция на производственное происшествие", "master@detailit.ru"),
    ("det-harness", "Обжим контактов: от инструмента до контроля усилия", "assembler@detailit.ru"),
    ("det-harness", "Чтение схем жгутов и предотвращение ошибок пиновки", "master@detailit.ru"),
    ("det-lab", "Функциональные испытания автомобильной электроники", "tester@detailit.ru"),
    ("det-lab", "Диагностика отказов на испытательном стенде", "engineer@detailit.ru"),
    ("det-qc", "Входной, операционный и приемочный контроль", "quality@detailit.ru"),
    ("det-qc", "8D, 5 Why и оформление корректирующих действий", "engineer@detailit.ru"),
    ("det-store", "Учет партий компонентов и защита от пересортицы", "storekeeper@detailit.ru"),
    ("det-store", "ESD-логистика электронных компонентов", "storekeeper@detailit.ru"),
    ("det-eng", "Разработка технологической карты операции", "engineer@detailit.ru"),
    ("det-eng", "Анализ технологических рисков PFMEA для участка", "engineer@detailit.ru"),
    ("det-smt", "SMT-процесс: паста, установка, оплавление, AOI", "smt@detailit.ru"),
    ("det-smt", "Дефекты SMT-монтажа и профилактика брака", "engineer@detailit.ru"),
    ("det-pcb", "Подготовка печатных плат к монтажу", "smt@detailit.ru"),
    ("det-pcb", "Очистка, хранение и маркировка печатных плат", "quality@detailit.ru"),
    ("det-calib", "Калибровка стендов и достоверность измерений", "metrolog@detailit.ru"),
    ("det-calib", "Методики измерений и неопределенность результата", "metrolog@detailit.ru"),
    ("det-repair", "Ремонт электронных модулей и повторный контроль", "repair@detailit.ru"),
    ("det-repair", "Безопасная доработка изделий после несоответствия", "engineer@detailit.ru"),
    ("det-pack", "Упаковка, маркировка и защита готовой продукции", "pack@detailit.ru"),
    ("det-pack", "Подготовка партии к отгрузке заказчику", "planner@detailit.ru"),
    ("det-log", "Внутренняя логистика без потери партии", "planner@detailit.ru"),
    ("det-log", "Канбан, WIP и движение полуфабрикатов", "linelead@detailit.ru"),
    ("det-plan", "Планирование смен и загрузка линий", "planner@detailit.ru"),
    ("det-plan", "Управление узкими местами производства", "planner@detailit.ru"),
    ("det-train", "Наставничество рабочих на производственной линии", "director@detailit.ru"),
    ("det-train", "Матрица компетенций и допуск к операциям", "engineer@detailit.ru"),
    ("det-maint", "Профилактика оборудования и журнал обслуживания", "maintenance@detailit.ru"),
    ("det-maint", "Диагностика отказов производственной оснастки", "maintenance@detailit.ru"),
    ("det-line1", "Смена на линии 1: запуск, контроль, выпуск", "linelead@detailit.ru"),
    ("det-line1", "Командная работа оператора и мастера линии 1", "master@detailit.ru"),
    ("det-line2", "Смена на линии 2: стабильность и качество", "linelead@detailit.ru"),
    ("det-line2", "Работа с отклонениями на линии 2", "quality@detailit.ru"),
    ("det-proto", "Опытная партия: от идеи до серийного маршрута", "engineer@detailit.ru"),
    ("det-proto", "Документирование прототипов и уроки запуска", "engineer@detailit.ru"),
]

for tag, title, owner in COURSE_THEMES:
    COURSES.append((tag, title, owner, [
        "Назначение и контекст операции",
        "Рабочее место, материалы и инструменты",
        "Пошаговый производственный процесс",
        "Контроль качества и типовые ошибки",
        "Практика, тест и допуск к работе",
    ]))


def q(value: str) -> str:
    return "'" + value.replace("'", "''") + "'"


def j(value: str) -> str:
    return value.replace("\\", "\\\\").replace('"', '\\"').replace("\n", "\\n")


def text_payload(course: str, module: str, index: int) -> str:
    return (
        '{"content":"# ' + j(module) + '\\n\\n'
        + 'Раздел курса «' + j(course) + '» разбирает производственную ситуацию ООО \\"ДетаЛит\\". '
        + 'Сотрудник изучает цель операции, границы ответственности, признаки отклонений и порядок записи результата.\\n\\n'
        + '**Ключевые действия:**\\n\\n'
        + '- сверить задание, маршрутную карту и актуальную версию документа;\\n'
        + '- проверить рабочее место, материалы, инструмент и средства защиты;\\n'
        + '- выполнить операцию по технологической последовательности;\\n'
        + '- остановить передачу изделия при спорном качестве;\\n'
        + '- зафиксировать результат так, чтобы партию можно было восстановить по истории.\\n\\n'
        + 'Практический итог раздела: сотрудник должен объяснить, что он делает, зачем это влияет на изделие автомобиля и какие данные нужно оставить после операции."}'
    )


def test_payload(course: str, module: str) -> str:
    return (
        '{"question":"Что нужно сделать, если в разделе «' + j(module) + '» обнаружено отклонение, влияющее на качество партии?",'
        '"isMultiSelect":false,"output":2,'
        '"options":[{"id":1,"text":"Передать изделие дальше и вернуться к вопросу позже","isCorrect":false},'
        '{"id":2,"text":"Остановить передачу, изолировать спорную единицу и сообщить ответственному","isCorrect":true},'
        '{"id":3,"text":"Исправить запись без уведомления мастера","isCorrect":false}]}'
    )


def file_payload(idx: int) -> str:
    filename, mime = FILES[idx % len(FILES)]
    return '{"filename":"' + filename + '","object_key":"course_files/detailit/' + filename + '","mime_type":"' + mime + '"}'


def video_payload(idx: int) -> str:
    return '{"videoUrl":"' + VIDEOS[idx % len(VIDEOS)] + '"}'


def course_description(title: str) -> str:
    return (
        f"Большой практический курс ООО \"ДетаЛит\" для производственных подразделений по профилю ОКВЭД 29.31. "
        f"Тема: {title}. В курсе минимум 5 модулей, учебные файлы, видеоразборы, контрольные вопросы и практические чек-листы для допуска сотрудника."
    )


def emit_slide(lines, module_var: str, idx: int, title: str, description: str, slide_type: str, payload: str):
    lines.append(
        f"    PERFORM pg_temp.detailit_slide({module_var}, {idx}, {q(title)}, {q(description)}, '{slide_type}', $json${payload}$json$::jsonb);"
    )


def main():
    lines = [
        "-- +goose Up",
        "-- +goose StatementBegin",
        "INSERT INTO accounts (email, password_hash, role, username)",
        "VALUES",
    ]
    account_rows = [
        f"    ({q(email)}, {q(PASSWORD_HASH)}, '{role}'::roles, {q(username)})"
        for email, username, role in ACCOUNTS
    ]
    lines.append(",\n".join(account_rows))
    lines.extend([
        "ON CONFLICT (email) DO UPDATE SET username = EXCLUDED.username, role = EXCLUDED.role;",
        "",
        "WITH departments (tag, title, description, owner_email) AS (",
        "    VALUES",
    ])
    dept_rows = [
        f"        ({q(tag)}, {q(title)}, {q(description)}, {q(owner)})"
        for tag, title, description, owner in DEPARTMENTS
    ]
    lines.append(",\n".join(dept_rows))
    lines.extend([
        ")",
        "INSERT INTO organizations (tag, title, description, owner_id)",
        "SELECT d.tag, d.title, d.description, a.id",
        "FROM departments d",
        "JOIN accounts a ON a.email = d.owner_email",
        "ON CONFLICT (tag) DO UPDATE SET title = EXCLUDED.title, description = EXCLUDED.description, owner_id = EXCLUDED.owner_id;",
        "",
        "WITH membership (tag, email) AS (",
        "    VALUES",
    ])
    dept_tags = ["detailit", "det-prod", "det-harness", "det-lab", "det-qc", "det-store", "det-eng", "det-safety"] + [d[0] for d in DEPARTMENTS]
    member_rows = []
    for tag in dept_tags:
        for email in [
            "admin@lms.dev",
            "director@detailit.ru",
            "master@detailit.ru",
            "engineer@detailit.ru",
            "quality@detailit.ru",
            "operator@detailit.ru",
            "assembler@detailit.ru",
            "linelead@detailit.ru",
        ]:
            member_rows.append(f"        ({q(tag)}, {q(email)})")
    lines.append(",\n".join(member_rows))
    lines.extend([
        ")",
        "INSERT INTO organization_accounts (organization_id, account_id)",
        "SELECT o.id, a.id FROM membership m JOIN organizations o ON o.tag = m.tag JOIN accounts a ON a.email = m.email",
        "ON CONFLICT DO NOTHING;",
        "",
        "CREATE OR REPLACE FUNCTION pg_temp.detailit_course(p_org_tag text, p_title text, p_description text, p_owner_email text)",
        "RETURNS bigint AS $$",
        "DECLARE v_course_id bigint; v_org_id bigint; v_owner_id bigint;",
        "BEGIN",
        "    SELECT id INTO v_org_id FROM organizations WHERE tag = p_org_tag;",
        "    SELECT id INTO v_owner_id FROM accounts WHERE email = p_owner_email;",
        "    SELECT id INTO v_course_id FROM courses WHERE title = p_title AND organization_id = v_org_id LIMIT 1;",
        "    IF v_course_id IS NULL THEN",
        "        INSERT INTO courses (title, description, owner, organization_id) VALUES (p_title, p_description, v_owner_id, v_org_id) RETURNING id INTO v_course_id;",
        "    ELSE",
        "        UPDATE courses SET description = p_description, owner = v_owner_id WHERE id = v_course_id;",
        "    END IF;",
        "    RETURN v_course_id;",
        "END;",
        "$$ LANGUAGE plpgsql;",
        "",
        "CREATE OR REPLACE FUNCTION pg_temp.detailit_module(p_course_id bigint, p_index int, p_title text)",
        "RETURNS bigint AS $$",
        "DECLARE v_module_id bigint;",
        "BEGIN",
        "    SELECT m.id INTO v_module_id FROM modules m JOIN course_modules cm ON cm.module_id = m.id WHERE cm.course_id = p_course_id AND m.title = p_title LIMIT 1;",
        "    IF v_module_id IS NULL THEN",
        "        INSERT INTO modules (title) VALUES (p_title) RETURNING id INTO v_module_id;",
        "        INSERT INTO course_modules (course_id, module_id, index) VALUES (p_course_id, v_module_id, p_index) ON CONFLICT DO NOTHING;",
        "    ELSE",
        "        UPDATE course_modules SET index = p_index WHERE course_id = p_course_id AND module_id = v_module_id;",
        "    END IF;",
        "    RETURN v_module_id;",
        "END;",
        "$$ LANGUAGE plpgsql;",
        "",
        "CREATE OR REPLACE FUNCTION pg_temp.detailit_slide(p_module_id bigint, p_index int, p_title text, p_description text, p_type slide_variation, p_payload jsonb)",
        "RETURNS bigint AS $$",
        "DECLARE v_slide_id bigint;",
        "BEGIN",
        "    SELECT s.id INTO v_slide_id FROM slides s JOIN module_slides ms ON ms.slide_id = s.id WHERE ms.module_id = p_module_id AND s.title = p_title LIMIT 1;",
        "    IF v_slide_id IS NULL THEN",
        "        INSERT INTO slides (title, description, slide_type, payload) VALUES (p_title, p_description, p_type, p_payload) RETURNING id INTO v_slide_id;",
        "        INSERT INTO module_slides (module_id, slide_id, index) VALUES (p_module_id, v_slide_id, p_index) ON CONFLICT DO NOTHING;",
        "    ELSE",
        "        UPDATE slides SET description = p_description, slide_type = p_type, payload = p_payload WHERE id = v_slide_id;",
        "        UPDATE module_slides SET index = p_index WHERE module_id = p_module_id AND slide_id = v_slide_id;",
        "    END IF;",
        "    RETURN v_slide_id;",
        "END;",
        "$$ LANGUAGE plpgsql;",
        "",
        "DO $detailit$",
        "DECLARE c_id bigint; m_id bigint;",
        "BEGIN",
    ])

    for course_idx, (tag, title, owner, modules) in enumerate(COURSES):
        lines.append(f"    c_id := pg_temp.detailit_course({q(tag)}, {q(title)}, {q(course_description(title))}, {q(owner)});")
        for module_idx, module in enumerate(modules, start=1):
            lines.append(f"    m_id := pg_temp.detailit_module(c_id, {module_idx}, {q(module)});")
            emit_slide(lines, "m_id", 1, f"{module}: цель и результат", f"Теория и практический смысл раздела курса {title}", "TEXT", text_payload(title, module, 1))
            emit_slide(lines, "m_id", 2, f"{module}: учебная схема или документ", "Файл для изучения, повторения или очного инструктажа", "FILE", file_payload(course_idx + module_idx))
            emit_slide(lines, "m_id", 3, f"{module}: видеоразбор", "Видеоматериал для самостоятельного изучения", "VIDEO_URL", video_payload(course_idx + module_idx))
            emit_slide(lines, "m_id", 4, f"{module}: чек-лист сотрудника", "Практический порядок действий на рабочем месте", "TEXT", text_payload(title, module, 4))
            emit_slide(lines, "m_id", 5, f"{module}: контрольный тест", "Проверка понимания ключевого правила", "TEST", test_payload(title, module))
            emit_slide(lines, "m_id", 6, f"{module}: типовые ошибки и допуск", "Закрепление темы перед переходом к следующему модулю", "TEXT", text_payload(title, module, 6))
        lines.append("")

    lines.extend([
        "END;",
        "$detailit$;",
        "",
        "SELECT setval('accounts_id_seq', (SELECT MAX(id) FROM accounts));",
        "SELECT setval('organizations_id_seq', (SELECT MAX(id) FROM organizations));",
        "SELECT setval('courses_id_seq', (SELECT MAX(id) FROM courses));",
        "SELECT setval('modules_id_seq', (SELECT MAX(id) FROM modules));",
        "SELECT setval('slides_id_seq', (SELECT MAX(id) FROM slides));",
        "-- +goose StatementEnd",
        "",
        "-- +goose Down",
        "-- +goose StatementBegin",
        "DELETE FROM courses WHERE title IN (",
    ])
    lines.append(",\n".join(f"    {q(title)}" for _, title, _, _ in COURSES[10:]))
    lines.extend([
        ");",
        "DELETE FROM organizations WHERE tag IN (",
        ",\n".join(f"    {q(tag)}" for tag, *_ in DEPARTMENTS),
        ");",
        "DELETE FROM accounts WHERE email IN (",
        ",\n".join(f"    {q(email)}" for email, *_ in ACCOUNTS),
        ");",
        "-- +goose StatementEnd",
        "",
    ])

    OUT.write_text("\n".join(lines), encoding="utf-8")
    print(f"Wrote {OUT}")


if __name__ == "__main__":
    main()
