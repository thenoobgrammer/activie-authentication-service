package session

import (
	"auth-service/internal/infra/redis"
	"auth-service/internal/token"
	"auth-service/internal/user"
	"auth-service/pkg/logs"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetTokenAndStartSession(params ClaimRequirementsParams) (*token.Token, error)
	GetAnonymousTokenAndStartSession() (*token.AnonymousToken, error)
	RefreshToken(refreshTokenStr string) (*token.Token, error)
	ValidateToken(tokenStr string) error
	InvalidateToken(tokenStr string) error
}

type service struct {
	repo      Repository
	userRepo  user.Repository
	tokenRepo token.Repository
}

func NewService(repo Repository, userRepo user.Repository, tokenRepo token.Repository) Service {
	return &service{
		repo:      repo,
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *service) GetTokenAndStartSession(params ClaimRequirementsParams) (*token.Token, error) {
	claims := NewClaims(ClaimsParams{
		UserID:    params.UserID,
		UserRoles: params.UserRoles,
		ClientID:  params.ClientID,
		Email:     params.Email,
	})

	accessToken, err := GenerateJWT(claims)
	if err != nil {
		logs.Error("StartSession", "failed.to.generate.token", err)
		return nil, fmt.Errorf("failed.to.generate.token: %w", err)
	}

	refreshToken, err := GenerateOpaqueToken()
	if err != nil {
		logs.Error("StartSession", "failed.to.generate.refresh.token", err)
		return nil, fmt.Errorf("failed.to.generate.refresh.token: %w", err)
	}

	if err := s.repo.CreateSession(NewSessionParams{
		SessionID:  uuid.New().String(),
		ClientID:   claims.ClientID,
		UserID:     claims.UserID,
		TokenJTI:   claims.ID,
		LastIP:     params.LastIP,
		DeviceType: params.DeviceType,
		StartTime:  time.Now(),
		Exp:        claims.ExpiresAt.Time,
		IsActive:   true,
	}); err != nil {
		return nil, fmt.Errorf("failed.to.add.session")
	}

	ok, err := s.tokenRepo.RevokeUserRefreshTokens(claims.UserID)
	if err != nil || !ok {
		logs.Error("", "failed.to.revoke.token", err)
		return nil, fmt.Errorf("failed.to.revoke.token")
	}

	if err := s.tokenRepo.AddRefreshToken(token.NewRefreshTokenRecordParams{
		UserID:    &claims.UserID,
		ClientID:  &claims.ClientID,
		TokenHash: refreshToken,
	}); err != nil {
		if err := s.repo.DeleteSessionByJTI(claims.ID); err != nil {
			logs.Error("", "failed.to.delete.session", err)
			return nil, fmt.Errorf("failed.to.delete.session")
		}
		logs.Error("", "failed.to.add.refresh.token", err)
		return nil, fmt.Errorf("failed.to.add.refresh.token")
	}

	return &token.Token{
		AccessToken:  *accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresIn:    int(claims.ExpiresAt.Time.Sub(time.Now()).Seconds()),
		Principal:    claims.Email,
	}, nil
}

func (s *service) GetAnonymousTokenAndStartSession() (*token.AnonymousToken, error) {
	anonID := uuid.New().String()

	claims := NewClaims(ClaimsParams{
		UserID:    anonID,
		UserRoles: []string{"anonymous"},
		ClientID:  "",
		Email:     "",
	})

	accessToken, err := GenerateJWT(claims)
	if err != nil {
		logs.Error("GetAnonymousToken", "failed.to.generate.token", err)
		return nil, fmt.Errorf("failed.to.generate.token: %w", err)
	}

	return &token.AnonymousToken{
		AccessToken: *accessToken,
		SessionID:   anonID,
		ExpiresIn:   int(claims.ExpiresAt.Time.Sub(time.Now()).Seconds()),
	}, nil
}

func (s *service) RefreshToken(tokenStr string) (*token.Token, error) {
	// 1. revoke token in db and get it back
	refreshToken, err := s.tokenRepo.RevokeRefreshToken(tokenStr)
	if err != nil || !refreshToken.IsRevoked {
		logs.Error("RefreshToken", "failed.to.revoke.token", err)
		return nil, fmt.Errorf("failed.to.revoke.token")
	}

	user, err := s.userRepo.PeekUserById(*refreshToken.UserID)
	if err != nil {
		logs.Error("RefreshToken", "failed.to.peek.user", err)
		return nil, fmt.Errorf("failed.to.peek.user")
	}

	// 2. generate new access_token
	claims := NewClaims(ClaimsParams{
		UserID:    user.ID,
		UserRoles: *user.SystemRoles,
		ClientID:  *refreshToken.ClientID,
		Email:     user.Email,
	})

	accessToken, err := GenerateJWT(claims)
	if err != nil {
		logs.Error("StartSession", "failed.to.generate.token", err)
		return nil, fmt.Errorf("failed.to.generate.token: %w", err)
	}

	newRefreshToken, err := GenerateOpaqueToken()
	if err != nil {
		logs.Error("StartSession", "failed.to.generate.refresh.token", err)
		return nil, fmt.Errorf("failed.to.generate.refresh.token: %w", err)
	}

	if err := s.tokenRepo.AddRefreshToken(token.NewRefreshTokenRecordParams{
		UserID:    &claims.UserID,
		ClientID:  &claims.ClientID,
		TokenHash: newRefreshToken,
	}); err != nil {
		if err := s.repo.DeleteSessionByJTI(claims.ID); err != nil {
			logs.Error("", "failed.to.delete.session", err)
			return nil, fmt.Errorf("failed.to.delete.session")
		}
		logs.Error("", "failed.to.add.refresh.token", err)
		return nil, fmt.Errorf("failed.to.add.refresh.token")
	}

	return &token.Token{
		AccessToken:  *accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "bearer",
		ExpiresIn:    int(claims.ExpiresAt.Time.Sub(time.Now()).Seconds()),
		Principal:    claims.Email,
	}, nil
}

func (s *service) InvalidateToken(tokenStr string) error {
	claims, err := ExtractClaims(tokenStr) // already validates the claims
	if err != nil {
		return fmt.Errorf("invalid.token")
	}

	if err := redis.Set("jti-blocklist", claims.ID, claims.ExpiresAt.Time.Sub(time.Now())); err != nil {
		return fmt.Errorf("failed.to.invalidate.token")
	}

	if err := s.repo.DeactivateSession(claims.ID); err != nil {
		return fmt.Errorf("cannot.deactivate.session")
	}

	return nil
}

func (s *service) ValidateToken(tokenStr string) error {
	claims, err := ExtractClaims(tokenStr) // already validates the claims
	if err != nil {
		return fmt.Errorf("invalid.token")
	}

	val, err := redis.GetString("jti-blocklist:" + claims.ID)
	if err == nil && val != "" {
		return fmt.Errorf("token is blocked")
	}

	return nil
}
