package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/model"
)

func intPtr(n int) *int        { return &n }
func uintPtr(n uint64) *uint64 { return &n }

var (
	existingCourse = func(ctx context.Context, id uint64) (model.Course, error) { return model.Course{ID: id}, nil }
	missingCourse  = func(ctx context.Context, id uint64) (model.Course, error) { return model.Course{}, apperr.ErrNotFound }
)

func TestProgressService_ListWordProgresses(t *testing.T) {
	t.Run("course_id と mastery（カンマ区切り）を絞り込み条件にして repository に渡す", func(t *testing.T) {
		var gotUser uint64
		var gotFilter model.ProgressFilter
		repo := &mockProgressRepository{
			findWordProgresses: func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error) {
				gotUser, gotFilter = userID, filter
				return nil, nil
			},
		}
		s := NewProgressService(repo, &mockCourseRepository{findByID: existingCourse}, &mockWordRepository{})

		got, err := s.ListWordProgresses(context.Background(), 5, ProgressQuery{CourseID: uintPtr(1), Mastery: "weak, unsure"})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		want := model.ProgressFilter{CourseID: uintPtr(1), Masteries: []model.Mastery{model.MasteryWeak, model.MasteryUnsure}}
		if gotUser != 5 || !reflect.DeepEqual(gotFilter, want) {
			t.Errorf("repository に渡した値 = user %d, filter %#v", gotUser, gotFilter)
		}
		if got == nil || len(got) != 0 {
			t.Errorf("0件なら空のスライスのはずが %#v", got)
		}
	})

	t.Run("mastery の指定がなければ絞り込まない", func(t *testing.T) {
		var gotFilter model.ProgressFilter
		repo := &mockProgressRepository{
			findWordProgresses: func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error) {
				gotFilter = filter
				return nil, nil
			},
		}
		s := NewProgressService(repo, &mockCourseRepository{}, &mockWordRepository{})

		if _, err := s.ListWordProgresses(context.Background(), 1, ProgressQuery{}); err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if gotFilter.CourseID != nil || gotFilter.Masteries != nil {
			t.Errorf("filter = %#v", gotFilter)
		}
	})

	t.Run("不正な mastery なら ValidationError で、repository は呼ばない", func(t *testing.T) {
		s := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{}, &mockWordRepository{})

		_, err := s.ListWordProgresses(context.Background(), 1, ProgressQuery{Mastery: "weak,good"})
		assertValidationError(t, err, `"good" は指定できません`)
	})

	t.Run("存在しないコースなら ErrNotFound", func(t *testing.T) {
		s := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{findByID: missingCourse}, &mockWordRepository{})

		_, err := s.ListWordProgresses(context.Background(), 1, ProgressQuery{CourseID: uintPtr(99)})
		if !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
	})
}

func TestProgressService_ListReviewWords_DefaultMastery(t *testing.T) {
	var gotFilter model.ProgressFilter
	repo := &mockProgressRepository{
		findReviewWords: func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.ReviewWord, error) {
			gotFilter = filter
			return nil, nil
		},
	}
	s := NewProgressService(repo, &mockCourseRepository{}, &mockWordRepository{})

	got, err := s.ListReviewWords(context.Background(), 1, ProgressQuery{})
	if err != nil {
		t.Fatalf("エラーが返った: %v", err)
	}
	if want := []model.Mastery{model.MasteryWeak, model.MasteryUnsure}; !reflect.DeepEqual(gotFilter.Masteries, want) {
		t.Errorf("mastery 未指定なら苦手と微妙のはずが %v", gotFilter.Masteries)
	}
	if got == nil {
		t.Error("0件なら空のスライスのはず")
	}
}

func TestProgressService_SaveWordProgress(t *testing.T) {
	t.Run("ユーザー・単語・習得度を repository に渡して保存する", func(t *testing.T) {
		var gotUser, gotWord uint64
		var gotMastery model.Mastery
		repo := &mockProgressRepository{
			upsertWordProgress: func(ctx context.Context, userID, wordID uint64, mastery model.Mastery) (model.WordProgress, error) {
				gotUser, gotWord, gotMastery = userID, wordID, mastery
				return model.WordProgress{UserID: userID, WordID: wordID, Mastery: mastery}, nil
			},
		}
		words := &mockWordRepository{
			findByID: func(ctx context.Context, id uint64) (model.Word, error) { return model.Word{ID: id}, nil },
		}
		s := NewProgressService(repo, &mockCourseRepository{}, words)

		got, err := s.SaveWordProgress(context.Background(), 1, 3, model.WordProgressInput{Mastery: "mastered"})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if gotUser != 1 || gotWord != 3 || gotMastery != model.MasteryMastered || got.Mastery != model.MasteryMastered {
			t.Errorf("渡した値 = %d %d %q, 返った値 = %#v", gotUser, gotWord, gotMastery, got)
		}
	})

	t.Run("mastery が不正なら ValidationError で、保存しない", func(t *testing.T) {
		s := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{}, &mockWordRepository{})

		_, err := s.SaveWordProgress(context.Background(), 1, 3, model.WordProgressInput{Mastery: "perfect"})
		assertValidationError(t, err, "mastery は weak / unsure / mastered")
	})

	t.Run("存在しない単語なら ErrNotFound で、保存しない", func(t *testing.T) {
		words := &mockWordRepository{
			findByID: func(ctx context.Context, id uint64) (model.Word, error) { return model.Word{}, apperr.ErrNotFound },
		}
		s := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{}, words)

		_, err := s.SaveWordProgress(context.Background(), 1, 999, model.WordProgressInput{Mastery: "weak"})
		if !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
	})
}

func TestProgressService_SaveCourseProgress(t *testing.T) {
	// コースの単語は3語
	threeWords := func(ctx context.Context, courseID uint64) ([]model.CourseWord, error) {
		return make([]model.CourseWord, 3), nil
	}

	t.Run("0〜単語数の範囲なら保存する（単語数と同じ＝最後まで終わった）", func(t *testing.T) {
		var gotPosition int
		repo := &mockProgressRepository{
			upsertCourseProgress: func(ctx context.Context, userID, courseID uint64, position int) (model.CourseProgress, error) {
				gotPosition = position
				return model.CourseProgress{UserID: userID, CourseID: courseID, CurrentPosition: position}, nil
			},
		}
		s := NewProgressService(repo, &mockCourseRepository{findByID: existingCourse, findWords: threeWords}, &mockWordRepository{})

		if _, err := s.SaveCourseProgress(context.Background(), 1, 1, model.CourseProgressInput{CurrentPosition: intPtr(3)}); err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if gotPosition != 3 {
			t.Errorf("保存した位置 = %d", gotPosition)
		}
	})

	tests := []struct {
		name       string
		position   *int
		wantDetail string
	}{
		{"current_position が無い", nil, "current_position は必須です"},
		{"マイナス", intPtr(-1), "0以上"},
		{"単語数を超える", intPtr(4), "単語数（3）以下"},
	}
	for _, tt := range tests {
		t.Run(tt.name+"なら ValidationError で、保存しない", func(t *testing.T) {
			s := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{findByID: existingCourse, findWords: threeWords}, &mockWordRepository{})

			_, err := s.SaveCourseProgress(context.Background(), 1, 1, model.CourseProgressInput{CurrentPosition: tt.position})
			assertValidationError(t, err, tt.wantDetail)
		})
	}

	t.Run("存在しないコースなら ErrNotFound", func(t *testing.T) {
		s := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{findByID: missingCourse}, &mockWordRepository{})

		_, err := s.SaveCourseProgress(context.Background(), 1, 99, model.CourseProgressInput{CurrentPosition: intPtr(0)})
		if !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
	})
}

func TestProgressService_Summary(t *testing.T) {
	updated := time.Date(2026, 10, 25, 9, 0, 0, 0, time.UTC)
	repo := &mockProgressRepository{
		findWordProgresses: func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error) {
			return []model.WordProgress{
				{WordID: 1, Mastery: model.MasteryMastered},
				{WordID: 2, Mastery: model.MasteryWeak},
				{WordID: 3, Mastery: model.MasteryWeak},
				{WordID: 4, Mastery: model.MasteryUnsure},
			}, nil
		},
		findCourseProgresses: func(ctx context.Context, userID uint64) ([]model.CourseProgress, error) {
			return []model.CourseProgress{
				{CourseID: 1, CurrentPosition: 2, UpdatedAt: updated},
				{CourseID: 2, CurrentPosition: 5, UpdatedAt: updated}, // 保存後に単語が減って3語になったコース
			}, nil
		},
	}
	courses := &mockCourseRepository{
		findAll: func(ctx context.Context) ([]model.Course, error) {
			return []model.Course{{ID: 1, Name: "動物"}, {ID: 2, Name: "果物"}, {ID: 3, Name: "乗り物"}}, nil
		},
		countWords: func(ctx context.Context) (map[uint64]int, error) {
			return map[uint64]int{1: 3, 2: 3}, nil // 乗り物は単語なし
		},
	}
	words := &mockWordRepository{count: func(ctx context.Context) (int, error) { return 6, nil }}
	s := NewProgressService(repo, courses, words)

	got, err := s.Summary(context.Background(), 1)
	if err != nil {
		t.Fatalf("エラーが返った: %v", err)
	}
	want := model.Summary{
		LearnedWords: 4, TotalWords: 6, Weak: 2, Unsure: 1, Mastered: 1,
		Courses: []model.CourseSummary{
			{CourseID: 1, Name: "動物", TotalWords: 3, CurrentPosition: 2, Status: model.CourseInProgress, UpdatedAt: &updated},
			{CourseID: 2, Name: "果物", TotalWords: 3, CurrentPosition: 3, Status: model.CourseCompleted, UpdatedAt: &updated},
			{CourseID: 3, Name: "乗り物", TotalWords: 0, CurrentPosition: 0, Status: model.CourseNotStarted},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}
