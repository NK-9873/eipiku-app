package repository

import (
	"context"
	"database/sql"

	"eipiku/learning-service/internal/model"
)

// CourseRepository は courses テーブルの読み書き。service はこのインターフェースだけに依存する
type CourseRepository interface {
	FindAll(ctx context.Context) ([]model.Course, error)
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
