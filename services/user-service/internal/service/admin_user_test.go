package service

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"eipiku/user-service/internal/apperr"
	"eipiku/user-service/internal/model"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Run("管理者は role を指定してユーザーを作成できる", func(t *testing.T) {
		var saved model.NewUser
		repo := &mockUserRepository{
			create: func(ctx context.Context, u model.NewUser) (model.User, error) {
				saved = u
				return model.User{ID: 9, Name: u.Name, Email: u.Email, Role: u.Role}, nil
			},
		}
		got, err := newTestService(repo).CreateUser(context.Background(), adminActor,
			model.AdminUserInput{Name: "佐藤", Email: "Sato@Example.com", Password: "password123", Role: "admin"})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if saved.Email != "sato@example.com" || saved.Role != model.RoleAdmin || got.Role != model.RoleAdmin {
			t.Errorf("saved %#v, got %#v", saved, got)
		}
		if bcrypt.CompareHashAndPassword([]byte(saved.PasswordHash), []byte("password123")) != nil {
			t.Error("パスワードが bcrypt でハッシュ化されていない")
		}
	})

	t.Run("role を省略したら一般ユーザー", func(t *testing.T) {
		var saved model.NewUser
		repo := &mockUserRepository{
			create: func(ctx context.Context, u model.NewUser) (model.User, error) { saved = u; return model.User{}, nil },
		}
		newTestService(repo).CreateUser(context.Background(), adminActor, model.AdminUserInput{Name: "佐藤", Email: "sato@example.com", Password: "password123"})
		if saved.Role != model.RoleUser {
			t.Errorf("role = %q", saved.Role)
		}
	})

	t.Run("入力チェックは新規登録と同じで、role も確かめる", func(t *testing.T) {
		_, err := newTestService(&mockUserRepository{}).CreateUser(context.Background(), adminActor,
			model.AdminUserInput{Name: "", Email: "bad", Password: "short", Role: "owner"})
		for _, d := range []string{"name は必須です", "email の形式", "8文字以上", "role は admin / user"} {
			assertValidationError(t, err, d)
		}
	})

	t.Run("メールアドレスの重複は ErrConflict", func(t *testing.T) {
		repo := &mockUserRepository{
			create: func(ctx context.Context, u model.NewUser) (model.User, error) {
				return model.User{}, apperr.ErrConflict
			},
		}
		_, err := newTestService(repo).CreateUser(context.Background(), adminActor, model.AdminUserInput{Name: "佐藤", Email: "sato@example.com", Password: "password123"})
		if !errors.Is(err, apperr.ErrConflict) {
			t.Errorf("ErrConflict が返るはずが %v", err)
		}
	})

	t.Run("一般ユーザーは ErrForbidden", func(t *testing.T) {
		_, err := newTestService(&mockUserRepository{}).CreateUser(context.Background(), userActor, model.AdminUserInput{Name: "佐藤", Email: "sato@example.com", Password: "password123"})
		if !errors.Is(err, apperr.ErrForbidden) {
			t.Errorf("ErrForbidden が返るはずが %v", err)
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	// targets の中から id で探して返す repository。Update で渡された値を記録する
	newRepo := func(admins int, targets ...model.User) (*mockUserRepository, *model.UserUpdate) {
		var updated model.UserUpdate
		find := func(ctx context.Context, id uint64) (model.User, error) {
			for _, u := range targets {
				if u.ID == id {
					return u, nil
				}
			}
			return model.User{}, apperr.ErrNotFound
		}
		return &mockUserRepository{
			findByID:    find,
			countAdmins: func(ctx context.Context) (int, error) { return admins, nil },
			update:      func(ctx context.Context, id uint64, u model.UserUpdate) error { updated = u; return nil },
		}, &updated
	}
	admin1 := model.User{ID: 1, Role: model.RoleAdmin}
	user2 := model.User{ID: 2, Role: model.RoleUser}

	t.Run("password が空ならパスワードは変えない", func(t *testing.T) {
		repo, updated := newRepo(1, user2)
		_, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 2,
			model.AdminUserInput{Name: " 花子 ", Email: "Hanako@Example.com", Role: "user"})
		if err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if updated.Name != "花子" || updated.Email != "hanako@example.com" || updated.Role != model.RoleUser || updated.PasswordHash != nil {
			t.Errorf("updated = %#v", updated)
		}
	})

	t.Run("password を入力したらハッシュ化して変える", func(t *testing.T) {
		repo, updated := newRepo(1, user2)
		if _, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 2,
			model.AdminUserInput{Name: "花子", Email: "hanako@example.com", Password: "new-password", Role: "user"}); err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if updated.PasswordHash == nil || bcrypt.CompareHashAndPassword([]byte(*updated.PasswordHash), []byte("new-password")) != nil {
			t.Errorf("パスワードが変わっていない: %#v", updated)
		}
	})

	t.Run("入力したパスワードが短ければ ValidationError", func(t *testing.T) {
		repo, _ := newRepo(1, user2)
		_, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 2,
			model.AdminUserInput{Name: "花子", Email: "hanako@example.com", Password: "short", Role: "user"})
		assertValidationError(t, err, "8文字以上")
	})

	t.Run("role が無ければ ValidationError", func(t *testing.T) {
		repo, _ := newRepo(1, user2)
		_, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 2,
			model.AdminUserInput{Name: "花子", Email: "hanako@example.com"})
		assertValidationError(t, err, "role は必須です")
	})

	t.Run("ただ1人の管理者が自分を user にするのは ValidationError で、更新しない", func(t *testing.T) {
		repo, _ := newRepo(1, admin1)
		repo.update = nil // 呼ばれたら panic
		_, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 1,
			model.AdminUserInput{Name: "管理者", Email: "admin@example.com", Role: "user"})
		assertValidationError(t, err, "管理者が1人もいなくなる")
	})

	t.Run("ほかに管理者がいれば、自分を user にできる", func(t *testing.T) {
		repo, updated := newRepo(2, admin1)
		if _, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 1,
			model.AdminUserInput{Name: "管理者", Email: "admin@example.com", Role: "user"}); err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
		if updated.Role != model.RoleUser {
			t.Errorf("updated = %#v", updated)
		}
	})

	t.Run("管理者のまま名前だけ変えるときは、管理者の数を数えない", func(t *testing.T) {
		repo, _ := newRepo(1, admin1)
		repo.countAdmins = nil // 呼ばれたら panic
		if _, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 1,
			model.AdminUserInput{Name: "管理者2", Email: "admin@example.com", Role: "admin"}); err != nil {
			t.Fatalf("エラーが返った: %v", err)
		}
	})

	t.Run("存在しないユーザーは ErrNotFound", func(t *testing.T) {
		repo, _ := newRepo(1)
		_, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 99,
			model.AdminUserInput{Name: "花子", Email: "hanako@example.com", Role: "user"})
		if !errors.Is(err, apperr.ErrNotFound) {
			t.Errorf("ErrNotFound が返るはずが %v", err)
		}
	})

	t.Run("メールアドレスが他のユーザーと重複したら ErrConflict", func(t *testing.T) {
		repo, _ := newRepo(1, user2)
		repo.update = func(ctx context.Context, id uint64, u model.UserUpdate) error { return apperr.ErrConflict }
		_, err := newTestService(repo).UpdateUser(context.Background(), adminActor, 2,
			model.AdminUserInput{Name: "花子", Email: "admin@example.com", Role: "user"})
		if !errors.Is(err, apperr.ErrConflict) {
			t.Errorf("ErrConflict が返るはずが %v", err)
		}
	})

	t.Run("一般ユーザーは ErrForbidden", func(t *testing.T) {
		_, err := newTestService(&mockUserRepository{}).UpdateUser(context.Background(), userActor, 2,
			model.AdminUserInput{Name: "花子", Email: "hanako@example.com", Role: "admin"})
		if !errors.Is(err, apperr.ErrForbidden) {
			t.Errorf("ErrForbidden が返るはずが %v", err)
		}
	})
}
