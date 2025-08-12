package models

import "github.com/beego/beego/v2/client/orm"

type RefreshTokens struct {
	ID        int    `orm:"column(id);auto"`
	UserID    int    `orm:"column(user_id);"`
	Token     string `orm:"column(token);unique"`
	ExpiresAt string `orm:"column(expires_at);"`
	CreatedAt string `orm:"column(created_at);default(current_timestamp)"`
	Revoked   bool   `orm:"column(revoked);default(false)"`
}

func init() {
	orm.RegisterModel(new(RefreshTokens))
}

// CreateRefreshTokens creates a new refresh token
func CreateRefreshTokens(token *RefreshTokens, o orm.Ormer) (int64, error) {
	id, err := o.Insert(token)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetRefreshTokensByUserID retrieves a refresh token by user ID
func GetRefreshTokens(token string, o orm.Ormer) (*RefreshTokens, error) {
	var tokenData RefreshTokens
	err := o.QueryTable(new(RefreshTokens)).Filter("token", token).One(&token)
	if err != nil {
		return nil, err
	}
	return &tokenData, nil
}

// RevokeRefreshTokens marks a refresh token as revoked
func RevokeRefreshTokens(token string, o orm.Ormer) error {
	_, err := o.QueryTable(new(RefreshTokens)).Filter("token", token).Update(orm.Params{
		"revoked": true,
	})
	if err != nil {
		return err
	}
	return nil
}

// DeleteRefreshTokens deletes a refresh token by ID
func DeleteRefreshTokens(tokenID int, o orm.Ormer) error {
	_, err := o.QueryTable(new(RefreshTokens)).Filter("id", tokenID).Delete()
	if err != nil {
		return err
	}
	return nil
}
