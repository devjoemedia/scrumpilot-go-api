package store

import (
	"context"
	"errors"
	"time"

	"github.com/devjoemedia/scrumpilot-go-api/models"
	"gorm.io/gorm"
)

type RefreshTokenStore struct {
	db *gorm.DB
}

func NewRefreshTokenStore(db *gorm.DB) (*RefreshTokenStore, error) {
	if err := db.AutoMigrate(&models.RefreshToken{}); err != nil {
		return nil, err
	}

	return &RefreshTokenStore{db: db}, nil
}

func (s *RefreshTokenStore) Create(ctx context.Context, token *models.RefreshToken) error {
	return s.db.WithContext(ctx).Create(token).Error
}

func (s *RefreshTokenStore) GetByID(ctx context.Context, id string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := s.db.WithContext(ctx).Where("id = ? AND revoked = ? AND expires_at > ?",
		id, false, time.Now()).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (s *RefreshTokenStore) Revoke(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (s *RefreshTokenStore) RevokeAllForUser(ctx context.Context, userID uint) error {
	return s.db.WithContext(ctx).
		Where("user_id = ? AND revoked = ?", userID, false).
		Update("revoked", true).Error
}
