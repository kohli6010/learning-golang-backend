//go:generate mockery --name=IRefreshTokensRepo --output=../mocks --outpkg=mocks --case=underscore

package repository

import "golang-crud-ddd/domain"

// IRefreshTokensRepo ...
type IRefreshTokensRepo interface {
	SaveRefreshTokens(data *domain.RefreshTokens) error
	GetRefreshTokens(token string) (*domain.RefreshTokens, error)
	RevokeRefreshTokens(token string) error
}
