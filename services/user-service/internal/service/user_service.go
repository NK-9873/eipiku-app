package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"eipiku/user-service/internal/apperr"
	"eipiku/user-service/internal/auth"
	"eipiku/user-service/internal/model"
	"eipiku/user-service/internal/repository"
)

// TokenIssuer はログイン成功時にトークンを発行する（auth.JWT が満たす）
type TokenIssuer interface {
	Issue(userID uint64, role model.Role) (token string, expiresAt time.Time, err error)
}

type UserService struct {
	repo       repository.UserRepository
	tokens     TokenIssuer
	bcryptCost int
	// dummyHash は存在しないメールアドレスでログインされたときも bcrypt の比較をして、
	// 応答時間の差からメールアドレスの登録有無がわからないようにするためのもの
	dummyHash []byte
}

func NewUserService(repo repository.UserRepository, tokens TokenIssuer) *UserService {
	return newUserService(repo, tokens, bcrypt.DefaultCost)
}

func newUserService(repo repository.UserRepository, tokens TokenIssuer, cost int) *UserService {
	dummy, err := bcrypt.GenerateFromPassword([]byte("dummy-password-for-timing"), cost)
	if err != nil {
		panic(err) // cost が範囲外のときだけ起きる
	}
	return &UserService{repo: repo, tokens: tokens, bcryptCost: cost, dummyHash: dummy}
}

// パスワードの長さ。bcrypt は 72 バイトより後ろを無視するので、それを超える入力は受け付けない
const (
	minPasswordLength = 8
	maxPasswordBytes  = 72
)

// Signup はユーザーを登録する。メールアドレスが登録済みなら apperr.ErrConflict
func (s *UserService) Signup(ctx context.Context, in model.SignupInput) (model.User, error) {
	v := &validator{}
	name := v.required("name", in.Name, 50)
	email := normalizeEmail(v.required("email", in.Email, 255))
	if email != "" && !validEmail(email) {
		v.add("email の形式が正しくありません")
	}
	switch {
	case len([]rune(in.Password)) < minPasswordLength:
		v.add("password は%d文字以上にしてください", minPasswordLength)
	case len(in.Password) > maxPasswordBytes:
		v.add("password は%dバイト以内にしてください", maxPasswordBytes)
	}
	if err := v.err(); err != nil {
		return model.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), s.bcryptCost)
	if err != nil {
		return model.User{}, fmt.Errorf("パスワードのハッシュ化: %w", err)
	}
	user, err := s.repo.Create(ctx, model.NewUser{Name: name, Email: email, PasswordHash: string(hash)})
	if err != nil {
		return model.User{}, fmt.Errorf("ユーザーの登録: %w", err)
	}
	return user, nil
}

// Login はメールアドレスとパスワードを確かめてトークンを発行する。
// どちらが違っても同じ apperr.ErrUnauthorized を返す
func (s *UserService) Login(ctx context.Context, in model.LoginInput) (model.LoginResult, error) {
	v := &validator{}
	email := normalizeEmail(v.required("email", in.Email, 255))
	if in.Password == "" {
		v.add("password は必須です")
	}
	if err := v.err(); err != nil {
		return model.LoginResult{}, err
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if errors.Is(err, apperr.ErrNotFound) {
		bcrypt.CompareHashAndPassword(s.dummyHash, []byte(in.Password))
		return model.LoginResult{}, apperr.ErrUnauthorized
	}
	if err != nil {
		return model.LoginResult{}, fmt.Errorf("ユーザーの取得: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return model.LoginResult{}, apperr.ErrUnauthorized
	}

	token, expiresAt, err := s.tokens.Issue(user.ID, user.Role)
	if err != nil {
		return model.LoginResult{}, fmt.Errorf("トークンの発行: %w", err)
	}
	return model.LoginResult{Token: token, TokenType: "Bearer", ExpiresAt: expiresAt, User: user}, nil
}

// Me はログイン中のユーザーの情報。トークンは有効でもユーザーが削除済みなら apperr.ErrNotFound
func (s *UserService) Me(ctx context.Context, actor auth.Claims) (model.User, error) {
	user, err := s.repo.FindByID(ctx, actor.UserID)
	if err != nil {
		return model.User{}, fmt.Errorf("ユーザーの取得: %w", err)
	}
	return user, nil
}

// UpdateSettings は画像タイプを更新し、更新後のユーザーを返す
func (s *UserService) UpdateSettings(ctx context.Context, actor auth.Claims, in model.SettingsInput) (model.User, error) {
	v := &validator{}
	if !in.ImageType.Valid() {
		v.add("image_type は photo / illustration のいずれかを指定してください")
	}
	if err := v.err(); err != nil {
		return model.User{}, err
	}
	if err := s.repo.UpdateImageType(ctx, actor.UserID, in.ImageType); err != nil {
		return model.User{}, fmt.Errorf("設定の更新: %w", err)
	}
	return s.Me(ctx, actor)
}

// ListUsers はユーザー一覧（管理者のみ。それ以外は apperr.ErrForbidden）
func (s *UserService) ListUsers(ctx context.Context, actor auth.Claims) ([]model.User, error) {
	if actor.Role != model.RoleAdmin {
		return nil, apperr.ErrForbidden
	}
	users, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザー一覧の取得: %w", err)
	}
	if users == nil {
		users = []model.User{}
	}
	return users, nil
}

// GetUser はユーザーの詳細（管理者のみ。それ以外は apperr.ErrForbidden）
func (s *UserService) GetUser(ctx context.Context, actor auth.Claims, id uint64) (model.User, error) {
	if actor.Role != model.RoleAdmin {
		return model.User{}, apperr.ErrForbidden
	}
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("ユーザーの取得: %w", err)
	}
	return user, nil
}

// メールアドレスは大文字・小文字を区別せずに扱うため、小文字にそろえて保存・検索する
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// validEmail は「名前 <a@b>」のような表記は受け付けず、アドレスだけの形式か確かめる
func validEmail(email string) bool {
	addr, err := mail.ParseAddress(email)
	return err == nil && addr.Address == email && strings.Contains(email[strings.LastIndex(email, "@"):], ".")
}
