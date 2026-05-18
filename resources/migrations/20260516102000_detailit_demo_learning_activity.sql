-- +goose Up
-- +goose StatementBegin
WITH demo_workers(email, worker_no) AS (
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
        cc.course_id,
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
completed_pairs(email, course_rank) AS (
    VALUES
        ('worker@detailit.ru', 1),
        ('worker@detailit.ru', 2),
        ('worker@detailit.ru', 3),
        ('operator@detailit.ru', 1),
        ('operator@detailit.ru', 4),
        ('operator@detailit.ru', 5),
        ('assembler@detailit.ru', 2),
        ('assembler@detailit.ru', 6),
        ('quality@detailit.ru', 3),
        ('quality@detailit.ru', 7),
        ('storekeeper@detailit.ru', 8),
        ('student@lms.dev', 1),
        ('tester@detailit.ru', 9),
        ('smt@detailit.ru', 10),
        ('repair@detailit.ru', 11),
        ('pack@detailit.ru', 12)
),
partial_pairs(email, course_rank, max_tests, wrong_tests) AS (
    VALUES
        ('worker@detailit.ru', 13, 3, 0),
        ('operator@detailit.ru', 14, 4, 1),
        ('assembler@detailit.ru', 15, 2, 1),
        ('quality@detailit.ru', 16, 4, 0),
        ('storekeeper@detailit.ru', 17, 2, 1),
        ('student@lms.dev', 18, 3, 2),
        ('tester@detailit.ru', 2, 2, 1),
        ('smt@detailit.ru', 3, 3, 0),
        ('repair@detailit.ru', 4, 2, 1),
        ('pack@detailit.ru', 5, 2, 0),
        ('maintenance@detailit.ru', 6, 3, 1)
),
completed_results AS (
    SELECT
        a.id AS account_id,
        ts.module_id,
        ts.slide_id,
        ts.correct_option_id AS selected_option_id,
        true AS is_right,
        1 + ((dw.worker_no + ts.test_no + cp.course_rank) % 2) AS attempts,
        now() - ((dw.worker_no + cp.course_rank + ts.test_no)::int * interval '7 hours') AS first_attempt_at,
        now() - ((dw.worker_no + cp.course_rank + ts.test_no)::int * interval '7 hours') + interval '18 minutes' AS last_attempt_at
    FROM completed_pairs cp
    JOIN demo_workers dw ON dw.email = cp.email
    JOIN accounts a ON a.email = dw.email
    JOIN test_slides ts ON ts.course_rank = cp.course_rank
),
partial_results AS (
    SELECT
        a.id AS account_id,
        ts.module_id,
        ts.slide_id,
        CASE
            WHEN ts.test_no <= pp.wrong_tests THEN
                CASE WHEN ts.correct_option_id = 1 THEN 2 ELSE 1 END
            ELSE ts.correct_option_id
        END AS selected_option_id,
        (ts.test_no > pp.wrong_tests) AS is_right,
        CASE WHEN ts.test_no <= pp.wrong_tests THEN 3 ELSE 1 + ((dw.worker_no + ts.test_no) % 2) END AS attempts,
        now() - ((dw.worker_no + pp.course_rank + ts.test_no)::int * interval '5 hours') AS first_attempt_at,
        now() - ((dw.worker_no + pp.course_rank + ts.test_no)::int * interval '5 hours') + interval '11 minutes' AS last_attempt_at
    FROM partial_pairs pp
    JOIN demo_workers dw ON dw.email = pp.email
    JOIN accounts a ON a.email = dw.email
    JOIN test_slides ts ON ts.course_rank = pp.course_rank AND ts.test_no <= pp.max_tests
),
demo_results AS (
    SELECT * FROM completed_results
    UNION ALL
    SELECT * FROM partial_results
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
FROM demo_results
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
USING accounts a
WHERE tr.account_id = a.id
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
  );
-- +goose StatementEnd
