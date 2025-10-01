package session

import (
	"auth-service/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	StartAuthenticatedSession(params SessionInfoParams) (*string, error)
	StartAnonymousSession(params AnonymousSessionInfoParams) (*string, error)
	GetActiveUserSession(sessionId string) *models.UserSession
	GetActiveAnonymousSession(sessionId string) *models.AnonymousSession
	EndActiveUserSession(sessionId string) bool
	EndActiveAnonymousSession(sessionId string) bool
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) StartAuthenticatedSession(params SessionInfoParams) (*string, error) {
	if params.UserID == "" {
		return nil, errors.New("userId.is.missing")
	}

	claims := NewUserClaims(UserClaimInfoParams{
		UserID:    params.UserID,
		UserRoles: params.UserRoles,
	})

	token := GenerateToken(claims)
	if token == nil {
		return nil, fmt.Errorf("failed.to.generate.token")
	}

	ok := s.repo.DeleteUserSessionFromUserId(params.UserID)
	if !ok {
		return nil, fmt.Errorf("failed.to.delete.previous.session")
	}

	sessionId := s.repo.CreateUserSession(SesssionParams{
		SessionID:  uuid.New().String(),
		UserID:     params.UserID,
		TokenJTI:   claims.ID,
		LastIP:     params.LastIP,
		DeviceType: params.DeviceType,
		StartTime:  time.Now(),
		Exp:        claims.ExpiresAt.Time,
		IsActive:   true,
	})
	if sessionId == nil {
		return nil, fmt.Errorf("failed.to.create.session")
	}

	activeSession := s.repo.GetUserSessionFromId(*sessionId)
	if activeSession == nil {
		return nil, fmt.Errorf("failed.to.retreive.active.session")
	}

	return token, nil
}

func (s *service) StartAnonymousSession(params AnonymousSessionInfoParams) (*string, error) {
	if params.ClientID == "" {
		return nil, errors.New("clientId.is.missing")
	}

	claims := NewAnonymousClaims(AnonymousClaimInfoParams{
		ClientID: params.ClientID,
	})

	token, err := GenerateAnonymousToken(claims)
	if token == nil || err != nil {
		return nil, fmt.Errorf("failed.to.generate.token")
	}

	sessionId := s.repo.CreateAnonymousSession(AnonymousSessionParams{
		SessionID:  uuid.New().String(),
		ClientID:   params.ClientID,
		TokenJTI:   claims.ID,
		LastIP:     params.LastIP,
		DeviceType: params.DeviceType,
		StartTime:  time.Now(),
		Exp:        claims.ExpiresAt.Time,
		IsActive:   true,
	})
	if sessionId == nil {
		return nil, fmt.Errorf("failed.to.create.anonymous.session")
	}

	return token, nil
}

func (s *service) GetActiveUserSession(sessionId string) *models.UserSession {
	return s.repo.GetUserSessionFromId(sessionId)
}

func (s *service) GetActiveAnonymousSession(sessionId string) *models.AnonymousSession {
	return s.repo.GetAnonymousSessionFromId(sessionId)
}

func (s *service) EndActiveUserSession(sessionId string) bool {
	return s.repo.DeleteUserSessionFromId(sessionId)
}

func (s *service) EndActiveAnonymousSession(sessionId string) bool {
	return s.repo.DeleteAnonymousSessionFromId(sessionId)
}
