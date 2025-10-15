package token

import (
	"fmt"
)

type Service interface {
	AddRefreshToken(params NewRefreshTokenRecordParams) error
	RevokeRefreshToken(tokenHash string) (*RefreshToken, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) AddRefreshToken(params NewRefreshTokenRecordParams) error {
	if err := s.repo.AddRefreshToken(params); err != nil {
		return fmt.Errorf("failed.to.add.refresh.token")
	}
	return nil
}

func (s *service) RevokeRefreshToken(tokenHash string) (*RefreshToken, error) {
	token, err := s.repo.RevokeRefreshToken(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed.to.revoke.refresh.token")
	}
	return token, nil
}

func (s *service) RevokeUserRefreshTokens(userId string) error {
	ok, err := s.repo.RevokeUserRefreshTokens(userId)
	if err != nil || !ok {
		return fmt.Errorf("failed.to.revoke.user.refresh.tokens")
	}
	return nil
}
