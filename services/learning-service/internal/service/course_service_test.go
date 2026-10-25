package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/model"
)

func TestCourseService_ListCourses(t *testing.T) {
	desc := "身近な動物の名前を覚えよう"
	tests := []struct {
		name string
		repo []model.Course
		want []model.Course
	}{
		{
			name: "repository のコースをそのまま返す",
			repo: []model.Course{{ID: 1, Name: "動物", Description: &desc}, {ID: 2, Name: "果物"}},
			want: []model.Course{{ID: 1, Name: "動物", Description: &desc}, {ID: 2, Name: "果物"}},
		},
		{
			name: "0件なら nil ではなく空のスライスを返す",
			repo: nil,
			want: []model.Course{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockCourseRepository{
				findAll: func(ctx context.Context) ([]model.Course, error) { return tt.repo, nil },
			}
			s := NewCourseService(repo, &mockWordRepository{})

			got, err := s.ListCourses(context.Background())
			if err != nil {
				t.Fatalf("エラーが返った: %v", err)
			}
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCourseService_CreateCourse(t *testing.T) {
	t.Run("前後の空白を除き、空の説明は nil にして repository に渡す", func(t *testing.T) {
		var saved model.CourseInput
		repo := &mockCourseRepository{
			create: func(ctx context.Context, in model.CourseInput) (model.Course, error) {
				saved = in
				return model.Course{ID: 3, Name: in.Name, Description: in.Description}, nil
			},
		}
		s := NewCourseService(repo, &mockWordRepository{})

		got, err := s.CreateCourse(context.Background(), adminActor, model.CourseInput{Name: "  乗り物  ", Description: ptr("  ")})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if saved.Name != "乗り物" || saved.Description != nil {
			t.Errorf("repository に渡した値 = %#v", saved)
		}
		if got.ID != 3 {
			t.Errorf("got %#v", got)
		}
	})

	tests := []struct {
		name       string
		in         model.CourseInput
		wantDetail string
	}{
		{"name が空", model.CourseInput{Name: "   "}, "name は必須です"},
		{"name が101文字", model.CourseInput{Name: strings.Repeat("あ", 101)}, "name は100文字以内"},
		{"description が256文字", model.CourseInput{Name: "乗り物", Description: ptr(strings.Repeat("a", 256))}, "description は255文字以内"},
	}
	for _, tt := range tests {
		t.Run(tt.name+"なら ValidationError で、repository は呼ばない", func(t *testing.T) {
			s := NewCourseService(&mockCourseRepository{}, &mockWordRepository{})

			_, err := s.CreateCourse(context.Background(), adminActor, tt.in)
			assertValidationError(t, err, tt.wantDetail)
		})
	}
}

func TestCourseService_UpdateCourse_NotFound(t *testing.T) {
	repo := &mockCourseRepository{
		update: func(ctx context.Context, id uint64, in model.CourseInput) (model.Course, error) {
			return model.Course{}, apperr.ErrNotFound
		},
	}
	s := NewCourseService(repo, &mockWordRepository{})

	_, err := s.UpdateCourse(context.Background(), adminActor, 99, model.CourseInput{Name: "乗り物"})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Errorf("ErrNotFound が返るはずが %v", err)
	}
}

func TestCourseService_DeleteCourse_NotFound(t *testing.T) {
	repo := &mockCourseRepository{
		delete: func(ctx context.Context, id uint64) error { return apperr.ErrNotFound },
	}
	s := NewCourseService(repo, &mockWordRepository{})

	if err := s.DeleteCourse(context.Background(), adminActor, 99); !errors.Is(err, apperr.ErrNotFound) {
		t.Errorf("ErrNotFound が返るはずが %v", err)
	}
}

func TestCourseService_ListCourseWords(t *testing.T) {
	t.Run("コースが無ければ ErrNotFound", func(t *testing.T) {
		repo := &mockCourseRepository{
			findByID: func(ctx context.Context, id uint64) (model.Course, error) { return model.Course{}, apperr.ErrNotFound },
		}
		s := NewCourseService(repo, &mockWordRepository{})

		if _, err := s.ListCourseWords(context.Background(), 99); !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
	})

	t.Run("単語が0件なら空のスライスを返す", func(t *testing.T) {
		repo := &mockCourseRepository{
			findByID:  func(ctx context.Context, id uint64) (model.Course, error) { return model.Course{ID: id}, nil },
			findWords: func(ctx context.Context, courseID uint64) ([]model.CourseWord, error) { return nil, nil },
		}
		s := NewCourseService(repo, &mockWordRepository{})

		got, err := s.ListCourseWords(context.Background(), 1)
		if err != nil || got == nil || len(got) != 0 {
			t.Errorf("got %#v, err %v", got, err)
		}
	})
}

func TestCourseService_ReplaceCourseWords(t *testing.T) {
	existingCourse := func(ctx context.Context, id uint64) (model.Course, error) { return model.Course{ID: id}, nil }
	allExist := func(ctx context.Context, ids []uint64) ([]uint64, error) { return ids, nil }

	t.Run("指定した順番のまま repository に渡し、置き換え後の一覧を返す", func(t *testing.T) {
		var replaced []uint64
		repo := &mockCourseRepository{
			findByID: existingCourse,
			replaceWords: func(ctx context.Context, courseID uint64, wordIDs []uint64) error {
				replaced = wordIDs
				return nil
			},
			findWords: func(ctx context.Context, courseID uint64) ([]model.CourseWord, error) {
				return []model.CourseWord{
					{DisplayOrder: 1, Word: model.Word{ID: 3}},
					{DisplayOrder: 2, Word: model.Word{ID: 1}},
				}, nil
			},
		}
		s := NewCourseService(repo, &mockWordRepository{findExistingIDs: allExist})

		got, err := s.ReplaceCourseWords(context.Background(), adminActor, 1, []uint64{3, 1})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if !reflect.DeepEqual(replaced, []uint64{3, 1}) {
			t.Errorf("repository に渡した順番 = %v", replaced)
		}
		if len(got) != 2 || got[0].ID != 3 {
			t.Errorf("got %#v", got)
		}
	})

	t.Run("空の配列なら単語をすべて外す", func(t *testing.T) {
		var replaced []uint64
		repo := &mockCourseRepository{
			findByID: existingCourse,
			replaceWords: func(ctx context.Context, courseID uint64, wordIDs []uint64) error {
				replaced = wordIDs
				return nil
			},
			findWords: func(ctx context.Context, courseID uint64) ([]model.CourseWord, error) { return nil, nil },
		}
		s := NewCourseService(repo, &mockWordRepository{findExistingIDs: allExist})

		got, err := s.ReplaceCourseWords(context.Background(), adminActor, 1, []uint64{})
		if err != nil || replaced == nil || len(replaced) != 0 || len(got) != 0 {
			t.Errorf("replaced %#v, got %#v, err %v", replaced, got, err)
		}
	})

	t.Run("word_ids が無ければ ValidationError", func(t *testing.T) {
		s := NewCourseService(&mockCourseRepository{}, &mockWordRepository{})

		_, err := s.ReplaceCourseWords(context.Background(), adminActor, 1, nil)
		assertValidationError(t, err, "word_ids は必須です")
	})

	t.Run("同じ単語が重複していたら ValidationError", func(t *testing.T) {
		s := NewCourseService(&mockCourseRepository{}, &mockWordRepository{})

		_, err := s.ReplaceCourseWords(context.Background(), adminActor, 1, []uint64{1, 2, 1})
		assertValidationError(t, err, "id=1）が重複")
	})

	t.Run("コースが無ければ ErrNotFound", func(t *testing.T) {
		repo := &mockCourseRepository{
			findByID: func(ctx context.Context, id uint64) (model.Course, error) { return model.Course{}, apperr.ErrNotFound },
		}
		s := NewCourseService(repo, &mockWordRepository{})

		if _, err := s.ReplaceCourseWords(context.Background(), adminActor, 99, []uint64{1}); !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
	})

	t.Run("存在しない単語があれば ValidationError で、置き換えない", func(t *testing.T) {
		repo := &mockCourseRepository{findByID: existingCourse}
		words := &mockWordRepository{
			findExistingIDs: func(ctx context.Context, ids []uint64) ([]uint64, error) { return []uint64{1}, nil },
		}
		s := NewCourseService(repo, words)

		_, err := s.ReplaceCourseWords(context.Background(), adminActor, 1, []uint64{1, 999})
		assertValidationError(t, err, "単語（id=999）は存在しません")
	})
}

// assertValidationError は err が ValidationError で、wantDetail を含む理由があるか確認する
func assertValidationError(t *testing.T, err error, wantDetail string) {
	t.Helper()
	var ve *apperr.ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("ValidationError が返るはずが %v", err)
	}
	for _, d := range ve.Details {
		if strings.Contains(d, wantDetail) {
			return
		}
	}
	t.Errorf("Details %q に %q が含まれていない", ve.Details, wantDetail)
}
