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

// Signup はユーザーを登録する（一般ユーザーとして）。メールアドレスが登録済みなら apperr.ErrConflict
func (s *UserService) Signup(ctx context.Context, in model.SignupInput) (model.User, error) {
	v := &validator{}
	name, email := validateProfile(v, in.Name, in.Email)
	validatePassword(v, in.Password)
	if err := v.err(); err != nil {
		return model.User{}, err
	}
	return s.create(ctx, name, email, in.Password, model.RoleUser)
}

// CreateUser は管理者がユーザーを作成する（管理者のみ）。入力チェックは新規登録と同じで、role も指定できる
func (s *UserService) CreateUser(ctx context.Context, actor auth.Claims, in model.AdminUserInput) (model.User, error) {
	if actor.Role != model.RoleAdmin {
		return model.User{}, apperr.ErrForbidden
	}
	if in.Role == "" {
		in.Role = model.RoleUser
	}
	v := &validator{}
	name, email := validateProfile(v, in.Name, in.Email)
	validatePassword(v, in.Password)
	validateRole(v, in.Role)
	if err := v.err(); err != nil {
		return model.User{}, err
	}
	return s.create(ctx, name, email, in.Password, in.Role)
}

func (s *UserService) create(ctx context.Context, name, email, password string, role model.Role) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return model.User{}, fmt.Errorf("パスワードのハッシュ化: %w", err)
	}
	user, err := s.repo.Create(ctx, model.NewUser{Name: name, Email: email, PasswordHash: string(hash), Role: role})
	if err != nil {
		return model.User{}, fmt.Errorf("ユーザーの登録: %w", err)
	}
	return user, nil
}

// UpdateUser は管理者がユーザーの名前・メール・role を変更する（管理者のみ）。
// password が空ならパスワードは変えない。管理者が1人もいなくなる変更は ValidationError
func (s *UserService) UpdateUser(ctx context.Context, actor auth.Claims, id uint64, in model.AdminUserInput) (model.User, error) {
	if actor.Role != model.RoleAdmin {
		return model.User{}, apperr.ErrForbidden
	}
	v := &validator{}
	name, email := validateProfile(v, in.Name, in.Email)
	if in.Password != "" {
		validatePassword(v, in.Password)
	}
	if in.Role == "" {
		v.add("role は必須です")
	} else {
		validateRole(v, in.Role)
	}
	if err := v.err(); err != nil {
		return model.User{}, err
	}

	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.User{}, fmt.Errorf("ユーザーの取得: %w", err)
	}
	// 管理者を一般ユーザーに変えるときは、ほかに管理者が残るか確かめる
	if target.Role == model.RoleAdmin && in.Role != model.RoleAdmin {
		admins, err := s.repo.CountAdmins(ctx)
		if err != nil {
			return model.User{}, fmt.Errorf("管理者の数の取得: %w", err)
		}
		if admins <= 1 {
			v.add("管理者が1人もいなくなるため、role を user に変更できません（先にほかのユーザーを管理者にしてください）")
			return model.User{}, v.err()
		}
	}

	update := model.UserUpdate{Name: name, Email: email, Role: in.Role}
	if in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), s.bcryptCost)
		if err != nil {
			return model.User{}, fmt.Errorf("パスワードのハッシュ化: %w", err)
		}
		h := string(hash)
		update.PasswordHash = &h
	}
	if err := s.repo.Update(ctx, id, update); err != nil {
		return model.User{}, fmt.Errorf("ユーザーの更新: %w", err)
	}
	return s.repo.FindByID(ctx, id)
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

// validateProfile は名前とメールアドレスを確かめ、整えた値を返す（新規登録・管理者の作成・更新で共通）
func validateProfile(v *validator, name, email string) (string, string) {
	name = v.required("name", name, 50)
	email = normalizeEmail(v.required("email", email, 255))
	if email != "" && !validEmail(email) {
		v.add("email の形式が正しくありません")
	}
	return name, email
}

func validatePassword(v *validator, password string) {
	switch {
	case len([]rune(password)) < minPasswordLength:
		v.add("password は%d文字以上にしてください", minPasswordLength)
	case len(password) > maxPasswordBytes:
		v.add("password は%dバイト以内にしてください", maxPasswordBytes)
	}
}

func validateRole(v *validator, role model.Role) {
	if !role.Valid() {
		v.add("role は admin / user のいずれかを指定してください")
	}
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
