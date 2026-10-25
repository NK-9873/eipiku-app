package service

import (
	"context"
	"fmt"
	"slices"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/auth"
	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/repository"
)

type CourseService struct {
	repo     repository.CourseRepository
	wordRepo repository.WordRepository
}

func NewCourseService(repo repository.CourseRepository, wordRepo repository.WordRepository) *CourseService {
	return &CourseService{repo: repo, wordRepo: wordRepo}
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

func (s *CourseService) GetCourse(ctx context.Context, id uint64) (model.Course, error) {
	course, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Course{}, fmt.Errorf("コースの取得: %w", err)
	}
	return course, nil
}

// CreateCourse / UpdateCourse / DeleteCourse / ReplaceCourseWords は管理者のみ（それ以外は apperr.ErrForbidden）
func (s *CourseService) CreateCourse(ctx context.Context, actor auth.Claims, in model.CourseInput) (model.Course, error) {
	if !actor.IsAdmin() {
		return model.Course{}, apperr.ErrForbidden
	}
	in, err := validateCourse(in)
	if err != nil {
		return model.Course{}, err
	}
	course, err := s.repo.Create(ctx, in)
	if err != nil {
		return model.Course{}, fmt.Errorf("コースの登録: %w", err)
	}
	return course, nil
}

func (s *CourseService) UpdateCourse(ctx context.Context, actor auth.Claims, id uint64, in model.CourseInput) (model.Course, error) {
	if !actor.IsAdmin() {
		return model.Course{}, apperr.ErrForbidden
	}
	in, err := validateCourse(in)
	if err != nil {
		return model.Course{}, err
	}
	course, err := s.repo.Update(ctx, id, in)
	if err != nil {
		return model.Course{}, fmt.Errorf("コースの更新: %w", err)
	}
	return course, nil
}

func (s *CourseService) DeleteCourse(ctx context.Context, actor auth.Claims, id uint64) error {
	if !actor.IsAdmin() {
		return apperr.ErrForbidden
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("コースの削除: %w", err)
	}
	return nil
}

// ListCourseWords はコースの単語を出題順に返す。コースが無ければ apperr.ErrNotFound
func (s *CourseService) ListCourseWords(ctx context.Context, courseID uint64) ([]model.CourseWord, error) {
	if _, err := s.repo.FindByID(ctx, courseID); err != nil {
		return nil, fmt.Errorf("コースの取得: %w", err)
	}
	words, err := s.repo.FindWords(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("コースの単語の取得: %w", err)
	}
	if words == nil {
		words = []model.CourseWord{}
	}
	return words, nil
}

// ReplaceCourseWords はコースの単語を wordIDs の並び順でまとめて置き換え、置き換え後の一覧を返す。
// 空の配列なら単語をすべて外す
func (s *CourseService) ReplaceCourseWords(ctx context.Context, actor auth.Claims, courseID uint64, wordIDs []uint64) ([]model.CourseWord, error) {
	if !actor.IsAdmin() {
		return nil, apperr.ErrForbidden
	}
	v := &validator{}
	if wordIDs == nil {
		v.add("word_ids は必須です（単語を外すときは空の配列を指定してください）")
	}
	seen := map[uint64]bool{}
	for _, id := range wordIDs {
		if seen[id] {
			v.add("word_ids に同じ単語（id=%d）が重複しています", id)
		}
		seen[id] = true
	}
	if err := v.err(); err != nil {
		return nil, err
	}

	if _, err := s.repo.FindByID(ctx, courseID); err != nil {
		return nil, fmt.Errorf("コースの取得: %w", err)
	}

	existing, err := s.wordRepo.FindExistingIDs(ctx, wordIDs)
	if err != nil {
		return nil, fmt.Errorf("単語の存在確認: %w", err)
	}
	for _, id := range wordIDs {
		if !slices.Contains(existing, id) {
			v.add("単語（id=%d）は存在しません", id)
		}
	}
	if err := v.err(); err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceWords(ctx, courseID, wordIDs); err != nil {
		return nil, fmt.Errorf("コースの単語の置き換え: %w", err)
	}
	return s.ListCourseWords(ctx, courseID)
}

// courses テーブルの VARCHAR の長さに合わせる
func validateCourse(in model.CourseInput) (model.CourseInput, error) {
	v := &validator{}
	in.Name = v.required("name", in.Name, 100)
	in.Description = v.optional("description", in.Description, 255)
	return in, v.err()
}
