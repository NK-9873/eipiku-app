package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"eipiku/user-service/internal/apperr"
	"eipiku/user-service/internal/auth"
	"eipiku/user-service/internal/model"
)

// mockUserRepository は repository.UserRepository のモック。
// 差し込んでいないメソッドが呼ばれたら panic してテストを落とす（呼ばれてはいけない場面の確認になる）
type mockUserRepository struct {
	create          func(ctx context.Context, u model.NewUser) (model.User, error)
	findByID        func(ctx context.Context, id uint64) (model.User, error)
	findByEmail     func(ctx context.Context, email string) (model.User, error)
	findAll         func(ctx context.Context) ([]model.User, error)
	updateImageType func(ctx context.Context, id uint64, imageType model.ImageType) error
}

func (m *mockUserRepository) Create(ctx context.Context, u model.NewUser) (model.User, error) {
	if m.create == nil {
		panic("UserRepository.Create が予期せず呼ばれた")
	}
	return m.create(ctx, u)
}

func (m *mockUserRepository) FindByID(ctx context.Context, id uint64) (model.User, error) {
	if m.findByID == nil {
		panic("UserRepository.FindByID が予期せず呼ばれた")
	}
	return m.findByID(ctx, id)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	if m.findByEmail == nil {
		panic("UserRepository.FindByEmail が予期せず呼ばれた")
	}
	return m.findByEmail(ctx, email)
}

func (m *mockUserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	if m.findAll == nil {
		panic("UserRepository.FindAll が予期せず呼ばれた")
	}
	return m.findAll(ctx)
}

func (m *mockUserRepository) UpdateImageType(ctx context.Context, id uint64, imageType model.ImageType) error {
	if m.updateImageType == nil {
		panic("UserRepository.UpdateImageType が予期せず呼ばれた")
	}
	return m.updateImageType(ctx, id, imageType)
}

// fakeTokens は渡された値が分かるトークンを返す
type fakeTokens struct{}

func (fakeTokens) Issue(userID uint64, role model.Role) (string, time.Time, error) {
	return "token-for-" + string(role), time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC), nil
}

// テストでは bcrypt のコストを最小にして速くする
func newTestService(repo *mockUserRepository) *UserService {
	return newUserService(repo, fakeTokens{}, bcrypt.MinCost)
}

var (
	userActor  = auth.Claims{UserID: 2, Role: model.RoleUser}
	adminActor = auth.Claims{UserID: 1, Role: model.RoleAdmin}
)

func TestUserService_Signup(t *testing.T) {
	t.Run("メールアドレスを小文字にそろえ、パスワードは bcrypt でハッシュ化して保存する", func(t *testing.T) {
		var saved model.NewUser
		repo := &mockUserRepository{
			create: func(ctx context.Context, u model.NewUser) (model.User, error) {
				saved = u
				return model.User{ID: 5, Name: u.Name, Email: u.Email}, nil
			},
		}

		got, err := newTestService(repo).Signup(context.Background(), model.SignupInput{
			Name: " 山田 花子 ", Email: " Hanako@Example.COM ", Password: "password123",
		})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if saved.Name != "山田 花子" || saved.Email != "hanako@example.com" {
			t.Errorf("repository に渡した値 = %#v", saved)
		}
		if saved.PasswordHash == "password123" || bcrypt.CompareHashAndPassword([]byte(saved.PasswordHash), []byte("password123")) != nil {
			t.Errorf("パスワードが bcrypt でハッシュ化されていない: %q", saved.PasswordHash)
		}
		if got.ID != 5 {
			t.Errorf("got %#v", got)
		}
	})

	tests := []struct {
		name       string
		in         model.SignupInput
		wantDetail string
	}{
		{"name が空", model.SignupInput{Email: "a@example.com", Password: "password123"}, "name は必須です"},
		{"email が空", model.SignupInput{Name: "花子", Password: "password123"}, "email は必須です"},
		{"email の形式が違う", model.SignupInput{Name: "花子", Email: "hanako@", Password: "password123"}, "email の形式"},
		{"email にドメインの . がない", model.SignupInput{Name: "花子", Email: "hanako@localhost", Password: "password123"}, "email の形式"},
		{"password が7文字", model.SignupInput{Name: "花子", Email: "a@example.com", Password: "1234567"}, "8文字以上"},
		{"password が73バイト", model.SignupInput{Name: "花子", Email: "a@example.com", Password: strings.Repeat("a", 73)}, "72バイト以内"},
	}
	for _, tt := range tests {
		t.Run(tt.name+"なら ValidationError で、登録しない", func(t *testing.T) {
			_, err := newTestService(&mockUserRepository{}).Signup(context.Background(), tt.in)
			assertValidationError(t, err, tt.wantDetail)
		})
	}

	t.Run("メールアドレスが登録済みなら ErrConflict", func(t *testing.T) {
		repo := &mockUserRepository{
			create: func(ctx context.Context, u model.NewUser) (model.User, error) {
				return model.User{}, apperr.ErrConflict
			},
		}
		_, err := newTestService(repo).Signup(context.Background(), model.SignupInput{Name: "花子", Email: "a@example.com", Password: "password123"})
		if !errors.Is(err, apperr.ErrConflict) {
			t.Errorf("ErrConflict が返るはずが %v", err)
		}
	})
}

func TestUserService_Login(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	stored := model.User{ID: 2, Email: "hanako@example.com", PasswordHash: string(hash), Role: model.RoleUser}
	repo := &mockUserRepository{
		findByEmail: func(ctx context.Context, email string) (model.User, error) {
			if email == stored.Email {
				return stored, nil
			}
			return model.User{}, apperr.ErrNotFound
		},
	}
	s := newTestService(repo)

	t.Run("正しいメールアドレス（大文字混じりでも可）とパスワードならトークンを返す", func(t *testing.T) {
		got, err := s.Login(context.Background(), model.LoginInput{Email: "Hanako@Example.com", Password: "password123"})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if got.Token != "token-for-user" || got.TokenType != "Bearer" || got.User.ID != 2 {
			t.Errorf("got %#v", got)
		}
	})

	t.Run("パスワード違いと、登録されていないメールアドレスは同じ ErrUnauthorized", func(t *testing.T) {
		for _, in := range []model.LoginInput{
			{Email: "hanako@example.com", Password: "wrong-password"},
			{Email: "nobody@example.com", Password: "password123"},
		} {
			if _, err := s.Login(context.Background(), in); !errors.Is(err, apperr.ErrUnauthorized) {
				t.Errorf("%s: ErrUnauthorized が返るはずが %v", in.Email, err)
			}
		}
	})

	t.Run("email・password が空なら ValidationError", func(t *testing.T) {
		_, err := s.Login(context.Background(), model.LoginInput{})
		assertValidationError(t, err, "email は必須です")
		assertValidationError(t, err, "password は必須です")
	})
}

func TestUserService_UpdateSettings(t *testing.T) {
	t.Run("ログイン中のユーザーの画像タイプを更新し、更新後の情報を返す", func(t *testing.T) {
		var gotID uint64
		var gotType model.ImageType
		repo := &mockUserRepository{
			updateImageType: func(ctx context.Context, id uint64, imageType model.ImageType) error {
				gotID, gotType = id, imageType
				return nil
			},
			findByID: func(ctx context.Context, id uint64) (model.User, error) {
				return model.User{ID: id, ImageType: gotType}, nil
			},
		}

		got, err := newTestService(repo).UpdateSettings(context.Background(), userActor, model.SettingsInput{ImageType: "illustration"})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if gotID != 2 || gotType != model.ImageTypeIllustration || got.ImageType != model.ImageTypeIllustration {
			t.Errorf("更新した値 = %d %q, 返った値 = %#v", gotID, gotType, got)
		}
	})

	t.Run("image_type が不正なら ValidationError で、更新しない", func(t *testing.T) {
		_, err := newTestService(&mockUserRepository{}).UpdateSettings(context.Background(), userActor, model.SettingsInput{ImageType: "video"})
		assertValidationError(t, err, "image_type は photo / illustration")
	})
}

func TestUserService_Me_NotFound(t *testing.T) {
	repo := &mockUserRepository{
		findByID: func(ctx context.Context, id uint64) (model.User, error) { return model.User{}, apperr.ErrNotFound },
	}
	if _, err := newTestService(repo).Me(context.Background(), userActor); !errors.Is(err, apperr.ErrNotFound) {
		t.Errorf("ErrNotFound が返るはずが %v", err)
	}
}

func TestUserService_AdminOnly(t *testing.T) {
	t.Run("一般ユーザーは ErrForbidden で、repository は呼ばない", func(t *testing.T) {
		s := newTestService(&mockUserRepository{})
		if _, err := s.ListUsers(context.Background(), userActor); !errors.Is(err, apperr.ErrForbidden) {
			t.Errorf("ListUsers: ErrForbidden が返るはずが %v", err)
		}
		if _, err := s.GetUser(context.Background(), userActor, 1); !errors.Is(err, apperr.ErrForbidden) {
			t.Errorf("GetUser: ErrForbidden が返るはずが %v", err)
		}
	})

	t.Run("管理者は一覧と詳細を取得できる（0件なら空のスライス）", func(t *testing.T) {
		repo := &mockUserRepository{
			findAll:  func(ctx context.Context) ([]model.User, error) { return nil, nil },
			findByID: func(ctx context.Context, id uint64) (model.User, error) { return model.User{ID: id}, nil },
		}
		s := newTestService(repo)

		users, err := s.ListUsers(context.Background(), adminActor)
		if err != nil || users == nil || len(users) != 0 {
			t.Errorf("ListUsers: got %#v, err %v", users, err)
		}
		user, err := s.GetUser(context.Background(), adminActor, 2)
		if err != nil || user.ID != 2 {
			t.Errorf("GetUser: got %#v, err %v", user, err)
		}
	})

	t.Run("管理者でも存在しないユーザーは ErrNotFound", func(t *testing.T) {
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id uint64) (model.User, error) { return model.User{}, apperr.ErrNotFound },
		}
		if _, err := newTestService(repo).GetUser(context.Background(), adminActor, 99); !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
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
