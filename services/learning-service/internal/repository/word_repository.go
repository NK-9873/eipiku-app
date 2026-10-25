package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/model"
)

// WordRepository は words テーブルの読み書き
type WordRepository interface {
	FindAll(ctx context.Context) ([]model.Word, error)
	// FindByID は見つからなければ apperr.ErrNotFound を返す
	FindByID(ctx context.Context, id uint64) (model.Word, error)
	// FindExistingIDs は ids のうち words テーブルにあるものだけを返す
	FindExistingIDs(ctx context.Context, ids []uint64) ([]uint64, error)
	Create(ctx context.Context, in model.WordInput) (model.Word, error)
	// Update / Delete は対象が無ければ apperr.ErrNotFound を返す
	Update(ctx context.Context, id uint64, in model.WordInput) (model.Word, error)
	Delete(ctx context.Context, id uint64) error
}

type mysqlWordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) WordRepository {
	return &mysqlWordRepository{db: db}
}

const wordColumns = "id, english, japanese, example, photo_path, illustration_path"

func scanWord(row interface{ Scan(...any) error }) (model.Word, error) {
	var w model.Word
	err := row.Scan(&w.ID, &w.English, &w.Japanese, &w.Example, &w.PhotoPath, &w.IllustrationPath)
	return w, err
}

func (r *mysqlWordRepository) FindAll(ctx context.Context) ([]model.Word, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+wordColumns+" FROM words ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var words []model.Word
	for rows.Next() {
		w, err := scanWord(rows)
		if err != nil {
			return nil, err
		}
		words = append(words, w)
	}
	return words, rows.Err()
}

func (r *mysqlWordRepository) FindByID(ctx context.Context, id uint64) (model.Word, error) {
	w, err := scanWord(r.db.QueryRowContext(ctx, "SELECT "+wordColumns+" FROM words WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.Word{}, apperr.ErrNotFound
	}
	return w, err
}

func (r *mysqlWordRepository) FindExistingIDs(ctx context.Context, ids []uint64) ([]uint64, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?, ", len(ids)), ", ")
	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	rows, err := r.db.QueryContext(ctx, "SELECT id FROM words WHERE id IN ("+placeholders+")", args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var existing []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		existing = append(existing, id)
	}
	return existing, rows.Err()
}

func (r *mysqlWordRepository) Create(ctx context.Context, in model.WordInput) (model.Word, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO words (english, japanese, example, photo_path, illustration_path) VALUES (?, ?, ?, ?, ?)",
		in.English, in.Japanese, in.Example, in.PhotoPath, in.IllustrationPath)
	if err != nil {
		return model.Word{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Word{}, err
	}
	return wordFromInput(uint64(id), in), nil
}

func (r *mysqlWordRepository) Update(ctx context.Context, id uint64, in model.WordInput) (model.Word, error) {
	res, err := r.db.ExecContext(ctx,
		"UPDATE words SET english = ?, japanese = ?, example = ?, photo_path = ?, illustration_path = ? WHERE id = ?",
		in.English, in.Japanese, in.Example, in.PhotoPath, in.IllustrationPath, id)
	if err != nil {
		return model.Word{}, err
	}
	if err := requireAffected(res); err != nil {
		return model.Word{}, err
	}
	return wordFromInput(id, in), nil
}

// course_words と word_progresses は外部キーの ON DELETE CASCADE で一緒に消える
func (r *mysqlWordRepository) Delete(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM words WHERE id = ?", id)
	if err != nil {
		return err
	}
	return requireAffected(res)
}

func wordFromInput(id uint64, in model.WordInput) model.Word {
	return model.Word{
		ID:               id,
		English:          in.English,
		Japanese:         in.Japanese,
		Example:          in.Example,
		PhotoPath:        in.PhotoPath,
		IllustrationPath: in.IllustrationPath,
	}
}
