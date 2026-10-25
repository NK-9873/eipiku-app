package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"

	"eipiku/user-service/internal/apperr"
	"eipiku/user-service/internal/model"
)

// UserRepository は users テーブルの読み書き。service はこのインターフェースだけに依存する
type UserRepository interface {
	// Create はメールアドレスが登録済みなら apperr.ErrConflict を返す
	Create(ctx context.Context, u model.NewUser) (model.User, error)
	// FindByID / FindByEmail は見つからなければ apperr.ErrNotFound を返す
	FindByID(ctx context.Context, id uint64) (model.User, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindAll(ctx context.Context) ([]model.User, error)
	// UpdateImageType は対象が無ければ apperr.ErrNotFound を返す
	UpdateImageType(ctx context.Context, id uint64, imageType model.ImageType) error
}

type mysqlUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &mysqlUserRepository{db: db}
}

const userColumns = "id, name, email, password_hash, image_type, role, created_at, updated_at"

// MySQL の重複エラー（UNIQUE 制約違反）
const errDuplicateEntry = 1062

func scanUser(row interface{ Scan(...any) error }) (model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.ImageType, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *mysqlUserRepository) Create(ctx context.Context, u model.NewUser) (model.User, error) {
	res, err := r.db.ExecContext(ctx,
		"INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)", u.Name, u.Email, u.PasswordHash)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == errDuplicateEntry {
		return model.User{}, apperr.ErrConflict
	}
	if err != nil {
		return model.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.User{}, err
	}
	// image_type・role・日時は DB の初期値が入るので読み直す
	return r.FindByID(ctx, uint64(id))
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id uint64) (model.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx, "SELECT "+userColumns+" FROM users WHERE id = ?", id))
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, apperr.ErrNotFound
	}
	return u, err
}

func (r *mysqlUserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx, "SELECT "+userColumns+" FROM users WHERE email = ?", email))
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, apperr.ErrNotFound
	}
	return u, err
}

func (r *mysqlUserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+userColumns+" FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *mysqlUserRepository) UpdateImageType(ctx context.Context, id uint64, imageType model.ImageType) error {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET image_type = ? WHERE id = ?", imageType, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
