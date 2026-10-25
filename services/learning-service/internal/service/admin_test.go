package service

import (
	"context"
	"errors"
	"testing"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/model"
)

// 登録・更新・削除と管理者用の参照は、一般ユーザーなら ErrForbidden で、repository は呼ばない
// （モックは関数を差し込んでいないので、呼ばれたら panic してテストが落ちる）
func TestAdminOnly_ForbiddenForUser(t *testing.T) {
	ctx := context.Background()
	courses := NewCourseService(&mockCourseRepository{}, &mockWordRepository{})
	words := NewWordService(&mockWordRepository{})
	progress := NewProgressService(&mockProgressRepository{}, &mockCourseRepository{}, &mockWordRepository{})

	calls := map[string]func() error{
		"CreateCourse": func() error {
			_, err := courses.CreateCourse(ctx, userActor, model.CourseInput{Name: "乗り物"})
			return err
		},
		"UpdateCourse": func() error {
			_, err := courses.UpdateCourse(ctx, userActor, 1, model.CourseInput{Name: "乗り物"})
			return err
		},
		"DeleteCourse": func() error { return courses.DeleteCourse(ctx, userActor, 1) },
		"ReplaceCourseWords": func() error {
			_, err := courses.ReplaceCourseWords(ctx, userActor, 1, []uint64{1})
			return err
		},
		"CreateWord": func() error {
			_, err := words.CreateWord(ctx, userActor, model.WordInput{English: "car", Japanese: "車"})
			return err
		},
		"UpdateWord": func() error {
			_, err := words.UpdateWord(ctx, userActor, 1, model.WordInput{English: "car", Japanese: "車"})
			return err
		},
		"DeleteWord": func() error { return words.DeleteWord(ctx, userActor, 1) },
		"UserProgress": func() error {
			_, err := progress.UserProgress(ctx, userActor, 3)
			return err
		},
	}
	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if err := call(); !errors.Is(err, apperr.ErrForbidden) {
				t.Errorf("ErrForbidden が返るはずが %v", err)
			}
		})
	}
}

func TestAdminOnly_AllowedForAdmin(t *testing.T) {
	t.Run("管理者はコースを削除できる", func(t *testing.T) {
		var deleted uint64
		repo := &mockCourseRepository{delete: func(ctx context.Context, id uint64) error { deleted = id; return nil }}
		if err := NewCourseService(repo, &mockWordRepository{}).DeleteCourse(context.Background(), adminActor, 3); err != nil || deleted != 3 {
			t.Errorf("deleted %d, err %v", deleted, err)
		}
	})

	t.Run("管理者は指定ユーザーの習得度と進み具合を取得できる（記録が無ければ空の一覧）", func(t *testing.T) {
		var gotUsers []uint64
		repo := &mockProgressRepository{
			findWordProgresses: func(ctx context.Context, userID uint64, filter model.ProgressFilter) ([]model.WordProgress, error) {
				gotUsers = append(gotUsers, userID)
				return []model.WordProgress{{UserID: userID, WordID: 1, Mastery: model.MasteryWeak}}, nil
			},
			findCourseProgresses: func(ctx context.Context, userID uint64) ([]model.CourseProgress, error) {
				gotUsers = append(gotUsers, userID)
				return nil, nil
			},
		}
		got, err := NewProgressService(repo, &mockCourseRepository{}, &mockWordRepository{}).UserProgress(context.Background(), adminActor, 3)
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		// 管理者自身（id=1）ではなく、指定したユーザー（id=3）の記録を取る
		if len(gotUsers) != 2 || gotUsers[0] != 3 || gotUsers[1] != 3 {
			t.Errorf("repository に渡したユーザー = %v", gotUsers)
		}
		if got.UserID != 3 || len(got.WordProgresses) != 1 || got.CourseProgresses == nil || len(got.CourseProgresses) != 0 {
			t.Errorf("got %#v", got)
		}
	})
}
