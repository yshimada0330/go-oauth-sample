package repository

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"gorm.io/gorm"
)

type ClientStorage struct {
	db *gorm.DB
}

func NewDBClientStore(db *gorm.DB) *ClientStorage {
	return &ClientStorage{
		db: db,
	}
}

func (s *ClientStorage) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	var client Client
	err := s.db.First(&client, "client_id = ?", id).Error
	if err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	// そのままgormのentityを返すと、oauth2.Clientのインターフェイスを満たさないので変換する
	return &models.Client{
		ID:     client.ClientId,
		Secret: client.Secret,
		Domain: client.Domain,
	}, nil
}

// oauth2.ClientStoreのインターフェイス上は必要ないがclientを作成するためのメソッドは必要
// oauth2.ClientStoreの実装とは切り離して考える話なので、どこにどう実装するかは設計時に考える
