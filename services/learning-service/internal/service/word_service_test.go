package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/model"
)

func TestWordService_ListWords_Empty(t *testing.T) {
	repo := &mockWordRepository{
		findAll: func(ctx context.Context) ([]model.Word, error) { return nil, nil },
	}
	got, err := NewWordService(repo).ListWords(context.Background())
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("got %#v, err %v", got, err)
	}
}

func TestWordService_CreateWord(t *testing.T) {
	t.Run("前後の空白を除き、空の任意項目は nil にして repository に渡す", func(t *testing.T) {
		var saved model.WordInput
		repo := &mockWordRepository{
			create: func(ctx context.Context, in model.WordInput) (model.Word, error) {
				saved = in
				return model.Word{ID: 7, English: in.English, Japanese: in.Japanese}, nil
			},
		}

		_, err := NewWordService(repo).CreateWord(context.Background(), adminActor, model.WordInput{
			English:   " lion ",
			Japanese:  "ライオン",
			Example:   ptr("The lion is strong."),
			PhotoPath: ptr(""),
		})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if saved.English != "lion" || saved.Example == nil || *saved.Example != "The lion is strong." {
			t.Errorf("repository に渡した値 = %#v", saved)
		}
		if saved.PhotoPath != nil || saved.IllustrationPath != nil {
			t.Errorf("空や未指定のパスは nil になるはず: %#v", saved)
		}
	})

	t.Run("english と japanese が空なら両方の理由を返し、repository は呼ばない", func(t *testing.T) {
		_, err := NewWordService(&mockWordRepository{}).CreateWord(context.Background(), adminActor, model.WordInput{})
		assertValidationError(t, err, "english は必須です")
		assertValidationError(t, err, "japanese は必須です")
	})
}

func TestWordService_UpdateWord_NotFound(t *testing.T) {
	repo := &mockWordRepository{
		update: func(ctx context.Context, id uint64, in model.WordInput) (model.Word, error) {
			return model.Word{}, apperr.ErrNotFound
		},
	}
	_, err := NewWordService(repo).UpdateWord(context.Background(), adminActor, 99, model.WordInput{English: "lion", Japanese: "ライオン"})
	if !errors.Is(err, apperr.ErrNotFound) {
		t.Errorf("ErrNotFound が返るはずが %v", err)
	}
}

func TestWordService_GetWord_NotFound(t *testing.T) {
	repo := &mockWordRepository{
		findByID: func(ctx context.Context, id uint64) (model.Word, error) { return model.Word{}, apperr.ErrNotFound },
	}
	if _, err := NewWordService(repo).GetWord(context.Background(), 99); !errors.Is(err, apperr.ErrNotFound) {
		t.Errorf("ErrNotFound が返るはずが %v", err)
	}
}

func TestWordService_CreateWord_ExampleJapanese(t *testing.T) {
	t.Run("例文の日本語訳も前後の空白を除いて保存し、空なら nil にする", func(t *testing.T) {
		var saved []model.WordInput
		repo := &mockWordRepository{
			create: func(ctx context.Context, in model.WordInput) (model.Word, error) {
				saved = append(saved, in)
				return model.Word{}, nil
			},
		}
		s := NewWordService(repo)
		s.CreateWord(context.Background(), adminActor, model.WordInput{English: "dog", Japanese: "犬", ExampleJapanese: ptr(" 私は犬を飼っています。 ")})
		s.CreateWord(context.Background(), adminActor, model.WordInput{English: "dog", Japanese: "犬", ExampleJapanese: ptr("")})
		if len(saved) != 2 || saved[0].ExampleJapanese == nil || *saved[0].ExampleJapanese != "私は犬を飼っています。" || saved[1].ExampleJapanese != nil {
			t.Errorf("saved = %#v", saved)
		}
	})

	t.Run("256文字なら ValidationError", func(t *testing.T) {
		_, err := NewWordService(&mockWordRepository{}).CreateWord(context.Background(), adminActor,
			model.WordInput{English: "dog", Japanese: "犬", ExampleJapanese: ptr(strings.Repeat("あ", 256))})
		assertValidationError(t, err, "example_japanese は255文字以内")
	})
}
