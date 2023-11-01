package repository

import (
	"context"
	"log"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AccessTokenStorage mysql token store
type AccessTokenStorage struct {
	db *gorm.DB
}

func NewDBTokenStore(db *gorm.DB) *AccessTokenStorage {
	return &AccessTokenStorage{
		db: db,
	}
}

// Create create and store the new token information
func (s *AccessTokenStorage) Create(ctx context.Context, info oauth2.TokenInfo) error {
	buf, _ := jsoniter.Marshal(info)
	accessToken := &AccessToken{
		Data: string(buf),
	}

	{
		accessToken.Token = info.GetAccess()
		accessToken.ExpiredAt = info.GetAccessCreateAt().Add(info.GetAccessExpiresIn()).Unix()
		accessToken.ClientId = info.GetClientID()
	}

	return s.db.Omit(clause.Associations).Create(accessToken).Error
}

// RemoveByCode delete the authorization code
func (s *AccessTokenStorage) RemoveByCode(ctx context.Context, code string) error {
	log.Panic("RemoveByCode does not support")
	return nil
}

// RemoveByAccess use the access token to delete the token information
func (s *AccessTokenStorage) RemoveByAccess(ctx context.Context, access string) error {
	var accessToken AccessToken
	err := s.db.First(&accessToken, "token = ?", access).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil
	}

	return s.db.Delete(&accessToken).Error
}

// RemoveByRefresh use the refresh token to delete the token information
func (s *AccessTokenStorage) RemoveByRefresh(ctx context.Context, refresh string) error {
	log.Panic("RemoveByRefresh does not support")
	return nil
}

func (s *AccessTokenStorage) toTokenInfo(data string) oauth2.TokenInfo {
	var tm models.Token
	_ = jsoniter.Unmarshal([]byte(data), &tm)
	return &tm
}

// GetByCode use the authorization code for token information data
func (s *AccessTokenStorage) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	log.Panic("GetByCode does not support")
	return nil, nil
}

// GetByAccess use the access token for token information data
func (s *AccessTokenStorage) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	if access == "" {
		return nil, nil
	}

	var accessToken AccessToken
	err := s.db.First(&accessToken, "token = ?", access).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return s.toTokenInfo(accessToken.Data), nil
}

// GetByRefresh use the refresh token for token information data
func (s *AccessTokenStorage) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	log.Panic("GetByRefresh does not support")
	return nil, nil
}
