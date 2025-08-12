package repository

import (
	"golang-crud-ddd/domain"
	models "golang-crud-ddd/model"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// UserRepo ...
type RefreshTokensRepo struct {
	o orm.Ormer
}

// NewUserRepo ... creates a new UserRepo instance
func NewRefreshTokensRepo() *RefreshTokensRepo {
	return &RefreshTokensRepo{
		o: orm.NewOrm(),
	}
}

// SaveRefreshTokens saves a refresh token for a user
func (repo *RefreshTokensRepo) SaveRefreshTokens(data *domain.RefreshTokens) error {
	RefreshTokens := &models.RefreshTokens{
		UserID:    data.UserID,
		Token:     data.Token,
		ExpiresAt: data.ExpiresAt.Format("2006-01-02 15:04:05"), // convert time.Time to string
		CreatedAt: data.CreatedAt.Format("2006-01-02 15:04:05"), // convert time.Time to string
		Revoked:   data.Revoked,
	}
	_, err := models.CreateRefreshTokens(RefreshTokens, repo.o)
	return err
}

// GetRefreshTokens retrieves a refresh token for a user
func (repo *RefreshTokensRepo) GetRefreshTokens(token string) (*domain.RefreshTokens, error) {
	tokenData, err := models.GetRefreshTokens(token, repo.o)
	if err != nil {
		return nil, err
	}
	if tokenData == nil {
		return nil, nil
	}
	// Convert *models.RefreshTokens to *domain.RefreshTokens
	expiresAt, err := time.Parse("2006-01-02 15:04:05", tokenData.ExpiresAt)
	if err != nil {
		return nil, err
	}
	createdAt, err := time.Parse("2006-01-02 15:04:05", tokenData.CreatedAt)
	if err != nil {
		return nil, err
	}
	domainToken := &domain.RefreshTokens{
		ID:        tokenData.ID,
		UserID:    tokenData.UserID,
		Token:     tokenData.Token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
		Revoked:   tokenData.Revoked,
	}
	return domainToken, nil
}

// RevokeRefreshTokensByTokenID revokes a refresh token by its ID
func (repo *RefreshTokensRepo) RevokeRefreshTokens(token string) error {
	err := models.RevokeRefreshTokens(token, repo.o)
	if err != nil {
		return err
	}
	return nil
}
