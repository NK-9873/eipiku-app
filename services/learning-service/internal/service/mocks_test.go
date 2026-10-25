package service

import (
	"context"

	"eipiku/learning-service/internal/auth"
	"eipiku/learning-service/internal/model"
)

// テストで使うログイン中のユーザー
var (
	adminActor = auth.Claims{UserID: 1, Role: model.RoleAdmin}
	userActor  = auth.Claims{UserID: 2, Role: model.RoleUser}
)

// repository のモック。テストごとに必要なメソッドだけ関数を差し込む。
// 差し込んでいないメソッドが呼ばれたら panic してテストを落とす（呼ばれてはいけない場面の確認になる）

type mockCourseRepository struct {
	findAll      func(ctx context.Context) ([]model.Course, error)
	findByID     func(ctx context.Context, id uint64) (model.Course, error)
	create       func(ctx context.Context, in model.CourseInput) (model.Course, error)
	update       func(ctx context.Context, id uint64, in model.CourseInput) (model.Course, error)
	delete       func(ctx context.Context, id uint64) error
	findWords    func(ctx context.Context, courseID uint64) ([]model.CourseWord, error)
	replaceWords func(ctx context.Context, courseID uint64, wordIDs []uint64) error
	countWords   func(ctx context.Context) (map[uint64]int, error)
}

func (m *mockCourseRepository) FindAll(ctx context.Context) ([]model.Course, error) {
	if m.findAll == nil {
		panic("CourseRepository.FindAll が予期せず呼ばれた")
	}
	return m.findAll(ctx)
}

func (m *mockCourseRepository) FindByID(ctx context.Context, id uint64) (model.Course, error) {
	if m.findByID == nil {
		panic("CourseRepository.FindByID が予期せず呼ばれた")
	}
	return m.findByID(ctx, id)
}

func (m *mockCourseRepository) Create(ctx context.Context, in model.CourseInput) (model.Course, error) {
	if m.create == nil {
		panic("CourseRepository.Create が予期せず呼ばれた")
	}
	return m.create(ctx, in)
}

func (m *mockCourseRepository) Update(ctx context.Context, id uint64, in model.CourseInput) (model.Course, error) {
	if m.update == nil {
		panic("CourseRepository.Update が予期せず呼ばれた")
	}
	return m.update(ctx, id, in)
}

func (m *mockCourseRepository) Delete(ctx context.Context, id uint64) error {
	if m.delete == nil {
		panic("CourseRepository.Delete が予期せず呼ばれた")
	}
	return m.delete(ctx, id)
}

func (m *mockCourseRepository) FindWords(ctx context.Context, courseID uint64) ([]model.CourseWord, error) {
	if m.findWords == nil {
		panic("CourseRepository.FindWords が予期せず呼ばれた")
	}
	return m.findWords(ctx, courseID)
}

func (m *mockCourseRepository) ReplaceWords(ctx context.Context, courseID uint64, wordIDs []uint64) error {
	if m.replaceWords == nil {
		panic("CourseRepository.ReplaceWords が予期せず呼ばれた")
	}
	return m.replaceWords(ctx, courseID, wordIDs)
}

func (m *mockCourseRepository) CountWords(ctx context.Context) (map[uint64]int, error) {
	if m.countWords == nil {
		panic("CourseRepository.CountWords が予期せず呼ばれた")
	}
	return m.countWords(ctx)
}

type mockWordRepository struct {
	findAll         func(ctx context.Context) ([]model.Word, error)
	count           func(ctx context.Context) (int, error)
	findByID        func(ctx context.Context, id uint64) (model.Word, error)
	findExistingIDs func(ctx context.Context, ids []uint64) ([]uint64, error)
	create          func(ctx context.Context, in model.WordInput) (model.Word, error)
	update          func(ctx context.Context, id uint64, in model.WordInput) (model.Word, error)
	delete          func(ctx context.Context, id uint64) error
}

func (m *mockWordRepository) FindAll(ctx context.Context) ([]model.Word, error) {
	if m.findAll == nil {
		panic("WordRepository.FindAll が予期せず呼ばれた")
	}
	return m.findAll(ctx)
}

func (m *mockWordRepository) Count(ctx context.Context) (int, error) {
	if m.count == nil {
		panic("WordRepository.Count が予期せず呼ばれた")
	}
	return m.count(ctx)
}

func (m *mockWordRepository) FindByID(ctx context.Context, id uint64) (model.Word, error) {
	if m.findByID == nil {
		panic("WordRepository.FindByID が予期せず呼ばれた")
	}
	return m.findByID(ctx, id)
}

func (m *mockWordRepository) FindExistingIDs(ctx context.Context, ids []uint64) ([]uint64, error) {
	if m.findExistingIDs == nil {
		panic("WordRepository.FindExistingIDs が予期せず呼ばれた")
	}
	return m.findExistingIDs(ctx, ids)
}

func (m *mockWordRepository) Create(ctx context.Context, in model.WordInput) (model.Word, error) {
	if m.create == nil {
		panic("WordRepository.Create が予期せず呼ばれた")
	}
	return m.create(ctx, in)
}

func (m *mockWordRepository) Update(ctx context.Context, id uint64, in model.WordInput) (model.Word, error) {
	if m.update == nil {
		panic("WordRepository.Update が予期せず呼ばれた")
	}
	return m.update(ctx, id, in)
}

func (m *mockWordRepository) Delete(ctx context.Context, id uint64) error {
	if m.delete == nil {
		panic("WordRepository.Delete が予期せず呼ばれた")
	}
	return m.delete(ctx, id)
}

type mockProgressRepository struct {
	findWordProgresses   func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error)
	findReviewWords      func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.ReviewWord, error)
	upsertWordProgress   func(ctx context.Context, userID, wordID uint64, mastery model.Mastery) (model.WordProgress, error)
	findCourseProgresses func(ctx context.Context, userID uint64) ([]model.CourseProgress, error)
	upsertCourseProgress func(ctx context.Context, userID, courseID uint64, position int) (model.CourseProgress, error)
}

func (m *mockProgressRepository) FindWordProgresses(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error) {
	if m.findWordProgresses == nil {
		panic("ProgressRepository.FindWordProgresses が予期せず呼ばれた")
	}
	return m.findWordProgresses(ctx, userID, filter)
}

func (m *mockProgressRepository) FindReviewWords(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.ReviewWord, error) {
	if m.findReviewWords == nil {
		panic("ProgressRepository.FindReviewWords が予期せず呼ばれた")
	}
	return m.findReviewWords(ctx, userID, filter)
}

func (m *mockProgressRepository) UpsertWordProgress(ctx context.Context, userID, wordID uint64, mastery model.Mastery) (model.WordProgress, error) {
	if m.upsertWordProgress == nil {
		panic("ProgressRepository.UpsertWordProgress が予期せず呼ばれた")
	}
	return m.upsertWordProgress(ctx, userID, wordID, mastery)
}

func (m *mockProgressRepository) FindCourseProgresses(ctx context.Context, userID uint64) ([]model.CourseProgress, error) {
	if m.findCourseProgresses == nil {
		panic("ProgressRepository.FindCourseProgresses が予期せず呼ばれた")
	}
	return m.findCourseProgresses(ctx, userID)
}

func (m *mockProgressRepository) UpsertCourseProgress(ctx context.Context, userID, courseID uint64, position int) (model.CourseProgress, error) {
	if m.upsertCourseProgress == nil {
		panic("ProgressRepository.UpsertCourseProgress が予期せず呼ばれた")
	}
	return m.upsertCourseProgress(ctx, userID, courseID, position)
}

func ptr(s string) *string { return &s }
