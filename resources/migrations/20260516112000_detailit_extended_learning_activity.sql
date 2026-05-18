-- +goose Up
-- +goose StatementBegin
WITH employees(email, employee_no) AS (
    VALUES
        ('worker@detailit.ru', 1),
        ('operator@detailit.ru', 2),
        ('assembler@detailit.ru', 3),
        ('quality@detailit.ru', 4),
        ('storekeeper@detailit.ru', 5),
        ('student@lms.dev', 6),
        ('tester@detailit.ru', 7),
        ('smt@detailit.ru', 8),
        ('repair@detailit.ru', 9),
        ('pack@detailit.ru', 10),
        ('maintenance@detailit.ru', 11)
),
course_catalog AS (
    SELECT
        c.id AS course_id,
        row_number() OVER (ORDER BY c.id) AS course_rank
    FROM courses c
    WHERE EXISTS (
        SELECT 1
        FROM course_modules cm
        JOIN module_slides ms ON ms.module_id = cm.module_id
        JOIN slides s ON s.id = ms.slide_id AND s.slide_type = 'TEST'
        WHERE cm.course_id = c.id
    )
),
test_slides AS (
    SELECT
        cc.course_rank,
        cm.module_id,
        s.id AS slide_id,
        CASE
            WHEN (s.payload ->> 'output') ~ '^[0-9]+$' THEN (s.payload ->> 'output')::bigint
            ELSE 2
        END AS correct_option_id,
        row_number() OVER (
            PARTITION BY cc.course_id
            ORDER BY cm.index, ms.index, s.id
        ) AS test_no
    FROM course_catalog cc
    JOIN course_modules cm ON cm.course_id = cc.course_id
    JOIN module_slides ms ON ms.module_id = cm.module_id
    JOIN slides s ON s.id = ms.slide_id AND s.slide_type = 'TEST'
),
complete_plan(email, course_rank) AS (
    VALUES
        ('worker@detailit.ru', 19),
        ('operator@detailit.ru', 20),
        ('assembler@detailit.ru', 21),
        ('quality@detailit.ru', 22),
        ('storekeeper@detailit.ru', 23),
        ('student@lms.dev', 24),
        ('tester@detailit.ru', 25),
        ('smt@detailit.ru', 26),
        ('repair@detailit.ru', 27),
        ('pack@detailit.ru', 28),
        ('maintenance@detailit.ru', 29),
        ('quality@detailit.ru', 30)
),
active_plan(email, course_rank, max_tests, wrong_tests, attempt_shift) AS (
    VALUES
        ('worker@detailit.ru', 31, 4, 1, 1),
        ('operator@detailit.ru', 32, 3, 0, 2),
        ('assembler@detailit.ru', 33, 5, 2, 2),
        ('quality@detailit.ru', 34, 4, 1, 3),
        ('storekeeper@detailit.ru', 35, 3, 1, 1),
        ('student@lms.dev', 36, 2, 1, 2),
        ('tester@detailit.ru', 37, 5, 1, 3),
        ('smt@detailit.ru', 38, 4, 0, 1),
        ('repair@detailit.ru', 39, 3, 2, 2),
        ('pack@detailit.ru', 40, 3, 0, 1),
        ('maintenance@detailit.ru', 41, 4, 1, 3),
        ('worker@detailit.ru', 42, 2, 0, 1),
        ('operator@detailit.ru', 43, 3, 1, 2),
        ('assembler@detailit.ru', 44, 2, 1, 1)
),
complete_results AS (
    SELECT
        a.id AS account_id,
        ts.module_id,
        ts.slide_id,
        ts.correct_option_id AS selected_option_id,
        true AS is_right,
        1 + ((e.employee_no + ts.test_no + cp.course_rank) % 3) AS attempts,
        now() - ((e.employee_no + cp.course_rank + ts.test_no)::int * interval '4 hours') AS first_attempt_at,
        now() - ((e.employee_no + cp.course_rank + ts.test_no)::int * interval '4 hours') + interval '24 minutes' AS last_attempt_at
    FROM complete_plan cp
    JOIN employees e ON e.email = cp.email
    JOIN accounts a ON a.email = e.email
    JOIN test_slides ts ON ts.course_rank = cp.course_rank
),
active_results AS (
    SELECT
        a.id AS account_id,
        ts.module_id,
        ts.slide_id,
        CASE
            WHEN ts.test_no <= ap.wrong_tests THEN CASE WHEN ts.correct_option_id = 1 THEN 2 ELSE 1 END
            ELSE ts.correct_option_id
        END AS selected_option_id,
        (ts.test_no > ap.wrong_tests) AS is_right,
        CASE WHEN ts.test_no <= ap.wrong_tests THEN 2 + ap.attempt_shift ELSE 1 + ((e.employee_no + ts.test_no + ap.attempt_shift) % 2) END AS attempts,
        now() - ((e.employee_no + ap.course_rank + ts.test_no)::int * interval '3 hours') AS first_attempt_at,
        now() - ((e.employee_no + ap.course_rank + ts.test_no)::int * interval '3 hours') + interval '16 minutes' AS last_attempt_at
    FROM active_plan ap
    JOIN employees e ON e.email = ap.email
    JOIN accounts a ON a.email = e.email
    JOIN test_slides ts ON ts.course_rank = ap.course_rank AND ts.test_no <= ap.max_tests
),
results AS (
    SELECT * FROM complete_results
    UNION ALL
    SELECT * FROM active_results
)
INSERT INTO test_results (
    account_id,
    module_id,
    slide_id,
    selected_option_id,
    is_right,
    attempts,
    first_attempt_at,
    last_attempt_at
)
SELECT
    account_id,
    module_id,
    slide_id,
    selected_option_id,
    is_right,
    attempts,
    first_attempt_at,
    last_attempt_at
FROM results
ON CONFLICT (account_id, module_id, slide_id) DO UPDATE
SET selected_option_id = EXCLUDED.selected_option_id,
    is_right = EXCLUDED.is_right,
    attempts = EXCLUDED.attempts,
    first_attempt_at = EXCLUDED.first_attempt_at,
    last_attempt_at = EXCLUDED.last_attempt_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM test_results tr
USING accounts a, course_modules cm, courses c
WHERE tr.account_id = a.id
  AND tr.module_id = cm.module_id
  AND cm.course_id = c.id
  AND a.email IN (
      'worker@detailit.ru',
      'operator@detailit.ru',
      'assembler@detailit.ru',
      'quality@detailit.ru',
      'storekeeper@detailit.ru',
      'student@lms.dev',
      'tester@detailit.ru',
      'smt@detailit.ru',
      'repair@detailit.ru',
      'pack@detailit.ru',
      'maintenance@detailit.ru'
  )
  AND c.id IN (
      SELECT course_id
      FROM (
          SELECT id AS course_id, row_number() OVER (ORDER BY id) AS course_rank
          FROM courses
      ) ranked_courses
      WHERE course_rank BETWEEN 19 AND 44
  );
-- +goose StatementEnd
