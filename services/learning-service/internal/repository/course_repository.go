package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/model"
)

// CourseRepository は courses / course_words テーブルの読み書き。service はこのインターフェースだけに依存する
type CourseRepository interface {
	FindAll(ctx context.Context) ([]model.Course, error)
	// FindByID は見つからなければ apperr.ErrNotFound を返す
	FindByID(ctx context.Context, id uint64) (model.Course, error)
	Create(ctx context.Context, in model.CourseInput) (model.Course, error)
	// Update / Delete は対象が無ければ apperr.ErrNotFound を返す
	Update(ctx context.Context, id uint64, in model.CourseInput) (model.Course, error)
	Delete(ctx context.Context, id uint64) error
	// FindWords はコースの単語を display_order 順に返す
	FindWords(ctx context.Context, courseID uint64) ([]model.CourseWord, error)
	// ReplaceWords はコースの単語を wordIDs の並び順でまとめて置き換える（トランザクション）
	ReplaceWords(ctx context.Context, courseID uint64, wordIDs []uint64) error
	// CountWords はコースごとの単語数を返す（単語が0件のコースは含まれない）
	CountWords(ctx context.Context) (map[uint64]int, error)
}

type mysqlCourseRepository struct {
	db *sql.DB
}

func NewCourseRepository(db *sql.DB) CourseRepository {
	return &mysqlCourseRepository{db: db}
}

func (r *mysqlCourseRepository) FindAll(ctx context.Context) ([]model.Course, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description FROM courses ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}

func (r *mysqlCourseRepository) FindByID(ctx context.Context, id uint64) (model.Course, error) {
	var c model.Course
	err := r.db.QueryRowContext(ctx, "SELECT id, name, description FROM courses WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Course{}, apperr.ErrNotFound
	}
	return c, err
}

func (r *mysqlCourseRepository) Create(ctx context.Context, in model.CourseInput) (model.Course, error) {
	res, err := r.db.ExecContext(ctx, "INSERT INTO courses (name, description) VALUES (?, ?)", in.Name, in.Description)
	if err != nil {
		return model.Course{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Course{}, err
	}
	return model.Course{ID: uint64(id), Name: in.Name, Description: in.Description}, nil
}

func (r *mysqlCourseRepository) Update(ctx context.Context, id uint64, in model.CourseInput) (model.Course, error) {
	res, err := r.db.ExecContext(ctx, "UPDATE courses SET name = ?, description = ? WHERE id = ?", in.Name, in.Description, id)
	if err != nil {
		return model.Course{}, err
	}
	if err := requireAffected(res); err != nil {
		return model.Course{}, err
	}
	return model.Course{ID: id, Name: in.Name, Description: in.Description}, nil
}

// course_words と course_progresses は外部キーの ON DELETE CASCADE で一緒に消える
func (r *mysqlCourseRepository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM courses WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireAffected(res)
}

func (r *mysqlCourseRepository) FindWords(ctx context.Context, courseID uint64) ([]model.CourseWord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cw.display_order, w.id, w.english, w.japanese, w.example, w.photo_path, w.illustration_path
		FROM course_words cw
		JOIN words w ON w.id = cw.word_id
		WHERE cw.course_id = ?
		ORDER BY cw.display_order`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []model.CourseWord
	for rows.Next() {
		var cw model.CourseWord
		if err := rows.Scan(&cw.DisplayOrder, &cw.ID, &cw.English, &cw.Japanese, &cw.Example, &cw.PhotoPath, &cw.IllustrationPath); err != nil {
			return nil, err
		}
		words = append(words, cw)
	}
	return words, rows.Err()
}

func (r *mysqlCourseRepository) ReplaceWords(ctx context.Context, courseID uint64, wordIDs []uint64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Commit 後に呼ばれても何もしない

	// 同じコースへの置き換えが同時に来ても順番に処理されるよう、コースの行をロックする
	var id uint64
	err = tx.QueryRowContext(ctx, "SELECT id FROM courses WHERE id = ? FOR UPDATE", courseID).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return apperr.ErrNotFound
	}
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM course_words WHERE course_id = ?", courseID); err != nil {
		return err
	}
	if len(wordIDs) > 0 {
		placeholders := make([]string, len(wordIDs))
		args := make([]any, 0, len(wordIDs)*3)
		for i, wordID := range wordIDs {
			placeholders[i] = "(?, ?, ?)"
			args = append(args, courseID, wordID, i+1)
		}
		query := "INSERT INTO course_words (course_id, word_id, display_order) VALUES " + strings.Join(placeholders, ", ")
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *mysqlCourseRepository) CountWords(ctx context.Context) (map[uint64]int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT course_id, COUNT(*) FROM course_words GROUP BY course_id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := map[uint64]int{}
	for rows.Next() {
		var courseID uint64
		var n int
		if err := rows.Scan(&courseID, &n); err != nil {
			return nil, err
		}
		counts[courseID] = n
	}
	return counts, rows.Err()
}

// requireAffected は UPDATE / DELETE の対象が0件なら apperr.ErrNotFound を返す
func requireAffected(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
