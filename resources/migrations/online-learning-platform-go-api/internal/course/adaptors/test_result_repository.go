package adaptors

import (
	"context"
	"errors"
	"time"

	"online-learning-platform-go-api/internal/course/usecase"

	"gorm.io/gorm"
)

type TestResultRepository struct {
	db *gorm.DB
}

func NewTestResultRepository(db *gorm.DB) *TestResultRepository {
	return &TestResultRepository{db: db}
}

type testResultRow struct {
	ID               uint64    `gorm:"column:id"`
	AccountID        uint64    `gorm:"column:account_id"`
	ModuleID         uint64    `gorm:"column:module_id"`
	SlideID          uint64    `gorm:"column:slide_id"`
	SelectedOptionID uint64    `gorm:"column:selected_option_id"`
	IsRight          bool      `gorm:"column:is_right"`
	Attempts         int       `gorm:"column:attempts"`
	FirstAttemptAt   time.Time `gorm:"column:first_attempt_at"`
	LastAttemptAt    time.Time `gorm:"column:last_attempt_at"`
}

func (r *TestResultRepository) Upsert(ctx context.Context, accountID, moduleID, slideID, selectedOptionID uint64, isRight bool) error {
	return r.db.WithContext(ctx).Exec(`
INSERT INTO test_results (account_id, module_id, slide_id, selected_option_id, is_right, attempts, first_attempt_at, last_attempt_at)
VALUES (?, ?, ?, ?, ?, 1, now(), now())
ON CONFLICT (account_id, module_id, slide_id)
DO UPDATE SET
  selected_option_id = EXCLUDED.selected_option_id,
  is_right = EXCLUDED.is_right,
  attempts = test_results.attempts + 1,
  last_attempt_at = now()
`, accountID, moduleID, slideID, selectedOptionID, isRight).Error
}

func (r *TestResultRepository) GetByAccountAndSlide(ctx context.Context, accountID, slideID uint64) (*usecase.TestResultRecord, error) {
	var row testResultRow
	err := r.db.WithContext(ctx).
		Raw(`
SELECT id, account_id, module_id, slide_id, selected_option_id, is_right, attempts, first_attempt_at, last_attempt_at
FROM test_results
WHERE account_id = ? AND slide_id = ?
LIMIT 1
`, accountID, slideID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &usecase.TestResultRecord{
		AccountID:        row.AccountID,
		ModuleID:         row.ModuleID,
		SlideID:          row.SlideID,
		SelectedOptionID: row.SelectedOptionID,
		IsRight:          row.IsRight,
		Attempts:         row.Attempts,
		FirstAttemptAt:   row.FirstAttemptAt,
		LastAttemptAt:    row.LastAttemptAt,
	}, nil
}

type testResultListRow struct {
	AccountID        uint64    `gorm:"column:account_id"`
	AccountEmail     string    `gorm:"column:account_email"`
	AccountUsername  string    `gorm:"column:account_username"`
	CourseID         uint64    `gorm:"column:course_id"`
	CourseTitle      string    `gorm:"column:course_title"`
	ModuleID         uint64    `gorm:"column:module_id"`
	ModuleTitle      string    `gorm:"column:module_title"`
	SlideID          uint64    `gorm:"column:slide_id"`
	SlideTitle       string    `gorm:"column:slide_title"`
	SelectedOptionID uint64    `gorm:"column:selected_option_id"`
	IsRight          bool      `gorm:"column:is_right"`
	Attempts         int       `gorm:"column:attempts"`
	FirstAttemptAt   time.Time `gorm:"column:first_attempt_at"`
	LastAttemptAt    time.Time `gorm:"column:last_attempt_at"`
}

func (r *TestResultRepository) List(ctx context.Context, orgID *uint64) ([]usecase.TestResultRecord, error) {
	baseSQL := `
SELECT
  tr.account_id,
  a.email AS account_email,
  a.username AS account_username,
  c.id AS course_id,
  c.title AS course_title,
  tr.module_id,
  m.title AS module_title,
  tr.slide_id,
  s.title AS slide_title,
  tr.selected_option_id,
  tr.is_right,
  tr.attempts,
  tr.first_attempt_at,
  tr.last_attempt_at
FROM test_results tr
JOIN accounts a ON a.id = tr.account_id
JOIN modules m ON m.id = tr.module_id
JOIN slides s ON s.id = tr.slide_id
JOIN course_modules cm ON cm.module_id = tr.module_id
JOIN courses c ON c.id = cm.course_id
`

	var rows []testResultListRow
	var err error
	if orgID != nil {
		err = r.db.WithContext(ctx).
			Raw(baseSQL+`
WHERE EXISTS (
  SELECT 1
  FROM course_organizations co
  WHERE co.course_id = c.id AND co.organization_id = ?
)
ORDER BY tr.last_attempt_at DESC, tr.account_id ASC
`, *orgID).
			Scan(&rows).Error
	} else {
		err = r.db.WithContext(ctx).
			Raw(baseSQL + `
ORDER BY tr.last_attempt_at DESC, tr.account_id ASC
`).
			Scan(&rows).Error
	}
	if err != nil {
		return nil, err
	}

	out := make([]usecase.TestResultRecord, 0, len(rows))
	for i := range rows {
		out = append(out, usecase.TestResultRecord{
			AccountID:        rows[i].AccountID,
			AccountEmail:     rows[i].AccountEmail,
			AccountUsername:  rows[i].AccountUsername,
			CourseID:         rows[i].CourseID,
			CourseTitle:      rows[i].CourseTitle,
			ModuleID:         rows[i].ModuleID,
			ModuleTitle:      rows[i].ModuleTitle,
			SlideID:          rows[i].SlideID,
			SlideTitle:       rows[i].SlideTitle,
			SelectedOptionID: rows[i].SelectedOptionID,
			IsRight:          rows[i].IsRight,
			Attempts:         rows[i].Attempts,
			FirstAttemptAt:   rows[i].FirstAttemptAt,
			LastAttemptAt:    rows[i].LastAttemptAt,
		})
	}
	return out, nil
}

type courseProgressListRow struct {
	AccountID       uint64    `gorm:"column:account_id"`
	AccountEmail    string    `gorm:"column:account_email"`
	AccountUsername string    `gorm:"column:account_username"`
	CourseID        uint64    `gorm:"column:course_id"`
	CourseTitle     string    `gorm:"column:course_title"`
	TotalTests      int       `gorm:"column:total_tests"`
	AttemptedTests  int       `gorm:"column:attempted_tests"`
	PassedTests     int       `gorm:"column:passed_tests"`
	Completed       bool      `gorm:"column:completed"`
	LastActivityAt  time.Time `gorm:"column:last_activity_at"`
}

func (r *TestResultRepository) ListCourseProgress(ctx context.Context, orgID *uint64) ([]usecase.CourseProgressRecord, error) {
	orgFilter := ""
	args := []interface{}{}
	if orgID != nil {
		orgFilter = `
WHERE EXISTS (
  SELECT 1
  FROM course_organizations co
  WHERE co.course_id = c.id AND co.organization_id = ?
)`
		args = append(args, *orgID)
	}

	query := `
WITH course_tests AS (
  SELECT
    c.id AS course_id,
    c.title AS course_title,
    COUNT(DISTINCT s.id) AS total_tests
  FROM courses c
  JOIN course_modules cm ON cm.course_id = c.id
  JOIN module_slides ms ON ms.module_id = cm.module_id
  JOIN slides s ON s.id = ms.slide_id AND s.slide_type = 'TEST'
` + orgFilter + `
  GROUP BY c.id, c.title
),
account_course AS (
  SELECT
    c.id AS course_id,
    tr.account_id,
    a.email AS account_email,
    a.username AS account_username,
    COUNT(DISTINCT tr.slide_id) AS attempted_tests,
    COUNT(DISTINCT tr.slide_id) FILTER (WHERE tr.is_right) AS passed_tests,
    MAX(tr.last_attempt_at) AS last_activity_at
  FROM test_results tr
  JOIN accounts a ON a.id = tr.account_id
  JOIN course_modules cm ON cm.module_id = tr.module_id
  JOIN courses c ON c.id = cm.course_id
  JOIN slides s ON s.id = tr.slide_id AND s.slide_type = 'TEST'
` + orgFilter + `
  GROUP BY c.id, tr.account_id, a.email, a.username
)
SELECT
  ac.account_id,
  ac.account_email,
  ac.account_username,
  ct.course_id,
  ct.course_title,
  ct.total_tests,
  ac.attempted_tests,
  ac.passed_tests,
  (ct.total_tests > 0 AND ac.passed_tests >= ct.total_tests) AS completed,
  ac.last_activity_at
FROM course_tests ct
JOIN account_course ac ON ac.course_id = ct.course_id
ORDER BY completed DESC, ac.last_activity_at DESC, ct.course_title ASC, ac.account_username ASC
`

	if orgID != nil {
		args = append(args, *orgID)
	}

	var rows []courseProgressListRow
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]usecase.CourseProgressRecord, 0, len(rows))
	for i := range rows {
		out = append(out, usecase.CourseProgressRecord{
			AccountID:       rows[i].AccountID,
			AccountEmail:    rows[i].AccountEmail,
			AccountUsername: rows[i].AccountUsername,
			CourseID:        rows[i].CourseID,
			CourseTitle:     rows[i].CourseTitle,
			TotalTests:      rows[i].TotalTests,
			AttemptedTests:  rows[i].AttemptedTests,
			PassedTests:     rows[i].PassedTests,
			Completed:       rows[i].Completed,
			LastActivityAt:  rows[i].LastActivityAt,
		})
	}
	return out, nil
}

var _ usecase.TestResultRepository = (*TestResultRepository)(nil)

var ErrNoResult = errors.New("test result not found")
