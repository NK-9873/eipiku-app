package service

import (
	"context"

	"eipiku/learning-service/internal/model"
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

type mockWordRepository struct {
	findAll         func(ctx context.Context) ([]model.Word, error)
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

func ptr(s string) *string { return &s }
