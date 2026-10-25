package service

import (
	"context"
	"fmt"
	"strings"

	"eipiku/learning-service/internal/apperr"
	"eipiku/learning-service/internal/auth"
	"eipiku/learning-service/internal/model"
	"eipiku/learning-service/internal/repository"
)

type ProgressService struct {
	repo       repository.ProgressRepository
	courseRepo repository.CourseRepository
	wordRepo   repository.WordRepository
}

func NewProgressService(repo repository.ProgressRepository, courseRepo repository.CourseRepository, wordRepo repository.WordRepository) *ProgressService {
	return &ProgressService{repo: repo, courseRepo: courseRepo, wordRepo: wordRepo}
}

// ProgressQuery は一覧の絞り込み条件（クエリ文字列のまま受け取り、service でチェックする）
type ProgressQuery struct {
	CourseID *uint64 // nil なら全コース
	Mastery  string  // "weak,unsure" のようなカンマ区切り。空なら指定なし
}

// ListWordProgresses は自分の習得度一覧。0件のときも空のスライスを返す
func (s *ProgressService) ListWordProgresses(ctx context.Context, userID uint64, q ProgressQuery) ([]model.WordProgress, error) {
	filter, err := s.buildFilter(ctx, q, nil)
	if err != nil {
		return nil, err
	}
	progresses, err := s.repo.FindWordProgresses(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("習得度の取得: %w", err)
	}
	if progresses == nil {
		progresses = []model.WordProgress{}
	}
	return progresses, nil
}

// ListReviewWords は復習対象の単語。mastery の指定がなければ苦手と微妙
func (s *ProgressService) ListReviewWords(ctx context.Context, userID uint64, q ProgressQuery) ([]model.ReviewWord, error) {
	filter, err := s.buildFilter(ctx, q, []model.Mastery{model.MasteryWeak, model.MasteryUnsure})
	if err != nil {
		return nil, err
	}
	words, err := s.repo.FindReviewWords(ctx, userID, filter)
	if err != nil {
		return nil, fmt.Errorf("復習対象の取得: %w", err)
	}
	if words == nil {
		words = []model.ReviewWord{}
	}
	return words, nil
}

// SaveWordProgress は習得度を登録・更新する。単語が無ければ apperr.ErrNotFound
func (s *ProgressService) SaveWordProgress(ctx context.Context, userID, wordID uint64, in model.WordProgressInput) (model.WordProgress, error) {
	v := &validator{}
	if !in.Mastery.Valid() {
		v.add("mastery は weak / unsure / mastered のいずれかを指定してください")
	}
	if err := v.err(); err != nil {
		return model.WordProgress{}, err
	}
	if _, err := s.wordRepo.FindByID(ctx, wordID); err != nil {
		return model.WordProgress{}, fmt.Errorf("単語の取得: %w", err)
	}
	progress, err := s.repo.UpsertWordProgress(ctx, userID, wordID, in.Mastery)
	if err != nil {
		return model.WordProgress{}, fmt.Errorf("習得度の保存: %w", err)
	}
	return progress, nil
}

// ListCourseProgresses は自分の進み具合一覧。0件のときも空のスライスを返す
func (s *ProgressService) ListCourseProgresses(ctx context.Context, userID uint64) ([]model.CourseProgress, error) {
	progresses, err := s.repo.FindCourseProgresses(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("進み具合の取得: %w", err)
	}
	if progresses == nil {
		progresses = []model.CourseProgress{}
	}
	return progresses, nil
}

// SaveCourseProgress は中断位置を登録・更新する。0〜コースの単語数の範囲のみ受け付ける
func (s *ProgressService) SaveCourseProgress(ctx context.Context, userID, courseID uint64, in model.CourseProgressInput) (model.CourseProgress, error) {
	v := &validator{}
	if in.CurrentPosition == nil {
		v.add("current_position は必須です")
	} else if *in.CurrentPosition < 0 {
		v.add("current_position は0以上を指定してください")
	}
	if err := v.err(); err != nil {
		return model.CourseProgress{}, err
	}

	if _, err := s.courseRepo.FindByID(ctx, courseID); err != nil {
		return model.CourseProgress{}, fmt.Errorf("コースの取得: %w", err)
	}
	words, err := s.courseRepo.FindWords(ctx, courseID)
	if err != nil {
		return model.CourseProgress{}, fmt.Errorf("コースの単語の取得: %w", err)
	}
	if *in.CurrentPosition > len(words) {
		v.add("current_position はコースの単語数（%d）以下を指定してください", len(words))
		return model.CourseProgress{}, v.err()
	}

	progress, err := s.repo.UpsertCourseProgress(ctx, userID, courseID, *in.CurrentPosition)
	if err != nil {
		return model.CourseProgress{}, fmt.Errorf("進み具合の保存: %w", err)
	}
	return progress, nil
}

// UserProgress は指定したユーザーの習得度と進み具合（管理者のみ。それ以外は apperr.ErrForbidden）。
// ユーザーは user-service 側にあるので存在確認はせず、記録が無ければ空の一覧を返す
func (s *ProgressService) UserProgress(ctx context.Context, actor auth.Claims, userID uint64) (model.UserProgress, error) {
	if !actor.IsAdmin() {
		return model.UserProgress{}, apperr.ErrForbidden
	}
	words, err := s.ListWordProgresses(ctx, userID, ProgressQuery{})
	if err != nil {
		return model.UserProgress{}, err
	}
	courses, err := s.ListCourseProgresses(ctx, userID)
	if err != nil {
		return model.UserProgress{}, err
	}
	return model.UserProgress{UserID: userID, WordProgresses: words, CourseProgresses: courses}, nil
}

// Summary はマイページ用の学習状況
func (s *ProgressService) Summary(ctx context.Context, userID uint64) (model.Summary, error) {
	wordProgresses, err := s.repo.FindWordProgresses(ctx, userID, model.ProgressFilter{})
	if err != nil {
		return model.Summary{}, fmt.Errorf("習得度の取得: %w", err)
	}
	courseProgresses, err := s.repo.FindCourseProgresses(ctx, userID)
	if err != nil {
		return model.Summary{}, fmt.Errorf("進み具合の取得: %w", err)
	}
	totalWords, err := s.wordRepo.Count(ctx)
	if err != nil {
		return model.Summary{}, fmt.Errorf("単語数の取得: %w", err)
	}
	courses, err := s.courseRepo.FindAll(ctx)
	if err != nil {
		return model.Summary{}, fmt.Errorf("コース一覧の取得: %w", err)
	}
	wordCounts, err := s.courseRepo.CountWords(ctx)
	if err != nil {
		return model.Summary{}, fmt.Errorf("コースの単語数の取得: %w", err)
	}

	summary := model.Summary{LearnedWords: len(wordProgresses), TotalWords: totalWords, Courses: []model.CourseSummary{}}
	for _, p := range wordProgresses {
		switch p.Mastery {
		case model.MasteryWeak:
			summary.Weak++
		case model.MasteryUnsure:
			summary.Unsure++
		case model.MasteryMastered:
			summary.Mastered++
		}
	}

	for _, c := range courses {
		cs := model.CourseSummary{CourseID: c.ID, Name: c.Name, TotalWords: wordCounts[c.ID], Status: model.CourseNotStarted}
		for _, p := range courseProgresses {
			if p.CourseID != c.ID {
				continue
			}
			updatedAt := p.UpdatedAt
			cs.UpdatedAt = &updatedAt
			// 保存後にコースの単語が減っていても、単語数を超えないようにする
			cs.CurrentPosition = min(p.CurrentPosition, cs.TotalWords)
			switch {
			case cs.CurrentPosition == 0:
				cs.Status = model.CourseNotStarted
			case cs.CurrentPosition >= cs.TotalWords:
				cs.Status = model.CourseCompleted
			default:
				cs.Status = model.CourseInProgress
			}
		}
		summary.Courses = append(summary.Courses, cs)
	}
	return summary, nil
}

// buildFilter はクエリをチェックして絞り込み条件にする。mastery の指定がなければ defaultMasteries を使う
func (s *ProgressService) buildFilter(ctx context.Context, q ProgressQuery, defaultMasteries []model.Mastery) (model.ProgressFilter, error) {
	v := &validator{}
	filter := model.ProgressFilter{CourseID: q.CourseID, Masteries: defaultMasteries}
	if strings.TrimSpace(q.Mastery) != "" {
		filter.Masteries = nil
		for _, part := range strings.Split(q.Mastery, ",") {
			m := model.Mastery(strings.TrimSpace(part))
			if !m.Valid() {
				v.add("mastery に指定できるのは weak / unsure / mastered です（%q は指定できません）", part)
				continue
			}
			filter.Masteries = append(filter.Masteries, m)
		}
	}
	if err := v.err(); err != nil {
		return model.ProgressFilter{}, err
	}
	if q.CourseID != nil {
		if _, err := s.courseRepo.FindByID(ctx, *q.CourseID); err != nil {
			return model.ProgressFilter{}, fmt.Errorf("コースの取得: %w", err)
		}
	}
	return filter, nil
}
