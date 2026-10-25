package service

import (
	"context"
	"fmt"

	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/repository"
)

type CourseService struct {
	repo repository.CourseRepository
}

func NewCourseService(repo repository.CourseRepository) *CourseService {
	return &CourseService{repo: repo}
}

// ListCourses はコース一覧を返す。0件のときも nil ではなく空のスライスを返す（JSON で [] にするため）
func (s *CourseService) ListCourses(ctx context.Context) ([]model.Course, error) {
	courses, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("コース一覧の取得: %w", err)
	}
	if courses == nil {
		courses = []model.Course{}
	}
	return courses, nil
}
