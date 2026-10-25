package repository

import (
	"context"
	"database/sql"
	"strings"

	"eipiku/learning-service/internal/model"
)

// ProgressRepository は word_progresses / course_progresses テーブルの読み書き
type ProgressRepository interface {
	// FindWordProgresses は習得度を word_id 順に返す
	FindWordProgresses(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error)
	// FindReviewWords は単語の情報つきで、コースの出題順に返す（全コースなら最初に出てくるコースの順）
	FindReviewWords(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.ReviewWord, error)
	// UpsertWordProgress はなければ作成、あれば更新して、保存後の行を返す
	UpsertWordProgress(ctx context.Context, userID, wordID uint64, mastery model.Mastery) (model.WordProgress, error)
	FindCourseProgresses(ctx context.Context, userID uint64) ([]model.CourseProgress, error)
	UpsertCourseProgress(ctx context.Context, userID, courseID uint64, position int) (model.CourseProgress, error)
}

type mysqlProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) ProgressRepository {
	return &mysqlProgressRepository{db: db}
}

// progressWhere は絞り込み条件の WHERE 句と引数を作る（word_progresses の別名は wp）
func progressWhere(userID uint64, filter model.ProgressFilter) (string, []any) {
	conds := []string{"wp.user_id = ?"}
	args := []any{userID}
	if filter.CourseID != nil {
		conds = append(conds, "EXISTS (SELECT 1 FROM course_words f WHERE f.word_id = wp.word_id AND f.course_id = ?)")
		args = append(args, *filter.CourseID)
	}
	if len(filter.Masteries) > 0 {
		conds = append(conds, "wp.mastery IN ("+strings.TrimSuffix(strings.Repeat("?, ", len(filter.Masteries)), ", ")+")")
		for _, m := range filter.Masteries {
			args = append(args, m)
		}
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

func (r *mysqlProgressRepository) FindWordProgresses(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error) {
	where, args := progressWhere(userID, filter)
	rows, err := r.db.QueryContext(ctx,
		"SELECT wp.user_id, wp.word_id, wp.mastery, wp.updated_at FROM word_progresses wp"+where+" ORDER BY wp.word_id", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []model.WordProgress
	for rows.Next() {
		var p model.WordProgress
		if err := rows.Scan(&p.UserID, &p.WordID, &p.Mastery, &p.UpdatedAt); err != nil {
			return nil, err
		}
		progresses = append(progresses, p)
	}
	return progresses, rows.Err()
}

func (r *mysqlProgressRepository) FindReviewWords(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.ReviewWord, error) {
	where, args := progressWhere(userID, filter)
	// 並び順に使うコース：指定があればそのコース、なければ単語が入っている中で id が一番小さいコース
	orderCourse := "(SELECT MIN(course_id) FROM course_words WHERE word_id = w.id)"
	joinArgs := []any{}
	if filter.CourseID != nil {
		orderCourse = "?"
		joinArgs = append(joinArgs, *filter.CourseID)
	}
	query := `
		SELECT w.id, w.english, w.japanese, w.example, w.example_japanese, w.photo_path, w.illustration_path, wp.mastery, wp.updated_at
		FROM word_progresses wp
		JOIN words w ON w.id = wp.word_id
		LEFT JOIN course_words cw ON cw.word_id = w.id AND cw.course_id = ` + orderCourse +
		where + `
		ORDER BY cw.course_id IS NULL, cw.course_id, cw.display_order, w.id`
	rows, err := r.db.QueryContext(ctx, query, append(joinArgs, args...)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []model.ReviewWord
	for rows.Next() {
		var rw model.ReviewWord
		if err := rows.Scan(&rw.ID, &rw.English, &rw.Japanese, &rw.Example, &rw.ExampleJapanese, &rw.PhotoPath, &rw.IllustrationPath, &rw.Mastery, &rw.UpdatedAt); err != nil {
			return nil, err
		}
		words = append(words, rw)
	}
	return words, rows.Err()
}

func (r *mysqlProgressRepository) UpsertWordProgress(ctx context.Context, userID, wordID uint64, mastery model.Mastery) (model.WordProgress, error) {
	// 同じ値で更新しても updated_at が変わるよう、明示的に CURRENT_TIMESTAMP を入れる
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO word_progresses (user_id, word_id, mastery) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE mastery = VALUES(mastery), updated_at = CURRENT_TIMESTAMP`,
		userID, wordID, mastery)
	if err != nil {
		return model.WordProgress{}, err
	}
	var p model.WordProgress
	err = r.db.QueryRowContext(ctx,
		"SELECT user_id, word_id, mastery, updated_at FROM word_progresses WHERE user_id = ? AND word_id = ?", userID, wordID).
		Scan(&p.UserID, &p.WordID, &p.Mastery, &p.UpdatedAt)
	return p, err
}

func (r *mysqlProgressRepository) FindCourseProgresses(ctx context.Context, userID uint64) ([]model.CourseProgress, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT user_id, course_id, current_position, updated_at FROM course_progresses WHERE user_id = ? ORDER BY course_id", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []model.CourseProgress
	for rows.Next() {
		var p model.CourseProgress
		if err := rows.Scan(&p.UserID, &p.CourseID, &p.CurrentPosition, &p.UpdatedAt); err != nil {
			return nil, err
		}
		progresses = append(progresses, p)
	}
	return progresses, rows.Err()
}

func (r *mysqlProgressRepository) UpsertCourseProgress(ctx context.Context, userID, courseID uint64, position int) (model.CourseProgress, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO course_progresses (user_id, course_id, current_position) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE current_position = VALUES(current_position), updated_at = CURRENT_TIMESTAMP`,
		userID, courseID, position)
	if err != nil {
		return model.CourseProgress{}, err
	}
	var p model.CourseProgress
	err = r.db.QueryRowContext(ctx,
		"SELECT user_id, course_id, current_position, updated_at FROM course_progresses WHERE user_id = ? AND course_id = ?", userID, courseID).
		Scan(&p.UserID, &p.CourseID, &p.CurrentPosition, &p.UpdatedAt)
	return p, err
}
