package service

import (
	"context"
	"fmt"

	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/repository"
)

type WordService struct {
	repo repository.WordRepository
}

func NewWordService(repo repository.WordRepository) *WordService {
	return &WordService{repo: repo}
}

// ListWords は単語一覧を返す。0件のときも空のスライスを返す
func (s *WordService) ListWords(ctx context.Context) ([]model.Word, error) {
	words, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("単語一覧の取得: %w", err)
	}
	if words == nil {
		words = []model.Word{}
	}
	return words, nil
}

func (s *WordService) GetWord(ctx context.Context, id uint64) (model.Word, error) {
	word, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Word{}, fmt.Errorf("単語の取得: %w", err)
	}
	return word, nil
}

func (s *WordService) CreateWord(ctx context.Context, in model.WordInput) (model.Word, error) {
	in, err := validateWord(in)
	if err != nil {
		return model.Word{}, err
	}
	word, err := s.repo.Create(ctx, in)
	if err != nil {
		return model.Word{}, fmt.Errorf("単語の登録: %w", err)
	}
	return word, nil
}

func (s *WordService) UpdateWord(ctx context.Context, id uint64, in model.WordInput) (model.Word, error) {
	in, err := validateWord(in)
	if err != nil {
		return model.Word{}, err
	}
	word, err := s.repo.Update(ctx, id, in)
	if err != nil {
		return model.Word{}, fmt.Errorf("単語の更新: %w", err)
	}
	return word, nil
}

func (s *WordService) DeleteWord(ctx context.Context, id uint64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("単語の削除: %w", err)
	}
	return nil
}

// words テーブルの VARCHAR の長さに合わせる
func validateWord(in model.WordInput) (model.WordInput, error) {
	v := &validator{}
	in.English = v.required("english", in.English, 100)
	in.Japanese = v.required("japanese", in.Japanese, 100)
	in.Example = v.optional("example", in.Example, 255)
	in.PhotoPath = v.optional("photo_path", in.PhotoPath, 255)
	in.IllustrationPath = v.optional("illustration_path", in.IllustrationPath, 255)
	return in, v.err()
}
