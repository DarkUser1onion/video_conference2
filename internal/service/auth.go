package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"video_conference/internal/config"
	"video_conference/internal/domain"
	"video_conference/internal/repository"
	"video_conference/pkg/jwt"
	"video_conference/pkg/logger"
)

type AuthService interface {
	Register(ctx context.Context, email, password, displayName string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)
	ValidateToken(ctx context.Context, tokenString string) (*domain.User, error)
	Logout(ctx context.Context, refreshToken string) error
}

type LoginResponse struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type authService struct {
	userRepo repository.UserRepository
	jwtCfg   config.JWTConfig
	log      logger.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtCfg config.JWTConfig, log logger.Logger) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
		log:      log,
	}
}

func (s *authService) Register(ctx context.Context, email, password, displayName string) (*domain.User, error) {
	// Проверка существования пользователя
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Хеширование пароля
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("Failed to hash password", "error", err)
		return nil, errors.New("failed to hash password")
	}

	// Создание пользователя
	user := &domain.User{
		ID:              uuid.New(),
		Email:           email,
		PasswordHash:    string(passwordHash),
		DisplayName:     displayName,
		GlobalRole:      domain.GlobalRoleUser,
		IsActive:        true,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error("Failed to create user", "error", err)
		return nil, errors.New("failed to create user")
	}

	// Убираем пароль из ответа
	user.PasswordHash = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	// Получение пользователя
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Проверка активности
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Генерация токенов
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.GlobalRole, s.jwtCfg.AccessSecret, s.jwtCfg.AccessTTL)
	if err != nil {
		s.log.Error("Failed to generate access token", "error", err)
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, s.jwtCfg.RefreshSecret, s.jwtCfg.RefreshTTL)
	if err != nil {
		s.log.Error("Failed to generate refresh token", "error", err)
		return nil, errors.New("failed to generate refresh token")
	}

	// Сохранение сессии
	tokenHash := hashToken(refreshToken)
	session := &domain.UserSession{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: tokenHash,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(s.jwtCfg.RefreshTTL),
	}

	if err := s.userRepo.CreateSession(ctx, session); err != nil {
		s.log.Error("Failed to create session", "error", err)
		return nil, errors.New("failed to create session")
	}

	// Обновление времени последнего входа
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Warn("Failed to update last login", "error", err)
	}

	user.PasswordHash = ""
	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Валидация refresh token
	claims, err := jwt.ValidateRefreshToken(refreshToken, s.jwtCfg.RefreshSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("invalid token subject")
	}

	// Проверка сессии в БД
	tokenHash := hashToken(refreshToken)
	session, err := s.userRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("session not found or expired")
	}

	// Получение пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Генерация новых токенов
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.GlobalRole, s.jwtCfg.AccessSecret, s.jwtCfg.AccessTTL)
	if err != nil {
		s.log.Error("Failed to generate access token", "error", err)
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(user.ID, s.jwtCfg.RefreshSecret, s.jwtCfg.RefreshTTL)
	if err != nil {
		s.log.Error("Failed to generate refresh token", "error", err)
		return nil, errors.New("failed to generate refresh token")
	}

	// Отзыв старой сессии и создание новой
	if err := s.userRepo.RevokeSession(ctx, session.ID, "refreshed"); err != nil {
		s.log.Warn("Failed to revoke old session", "error", err)
	}

	newTokenHash := hashToken(newRefreshToken)
	newSession := &domain.UserSession{
		ID:               uuid.New(),
		UserID:           user.ID,
		RefreshTokenHash: newTokenHash,
		CreatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(s.jwtCfg.RefreshTTL),
	}

	if err := s.userRepo.CreateSession(ctx, newSession); err != nil {
		s.log.Error("Failed to create new session", "error", err)
		return nil, errors.New("failed to create new session")
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := jwt.ValidateToken(tokenString, s.jwtCfg.AccessSecret)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	return user, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	session, err := s.userRepo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		return errors.New("session not found")
	}

	return s.userRepo.RevokeSession(ctx, session.ID, "logout")
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

