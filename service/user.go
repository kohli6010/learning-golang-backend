package service

import (
	"errors"
	"golang-crud-ddd/apprequest"
	"golang-crud-ddd/appresponse"
	"golang-crud-ddd/domain"
	"golang-crud-ddd/helper"
	"golang-crud-ddd/repository"
	"log"
	"time"
)

// UserService ...
type UserService interface {
	CreateUser(user *apprequest.UserRequest) (*appresponse.AuthenticateUserResponse, error)
	AuthenticateUser(auth *apprequest.AuthenticateUserRequest) (*appresponse.AuthenticateUserResponse, error)
	ChangePassword(in *apprequest.ChangePasswordRequest) error
	RefreshTokenss(oldRefreshTokens string) (*appresponse.AuthenticateUserResponse, error)
	Logout(RefreshTokens string) error
	GetUserByID(id int) (*appresponse.UserResponse, error)
	UpdateUserByID(id int, user *apprequest.UpsertUserPersonalDetailRequest) (*appresponse.UserResponse, error)
}

// userService
type userService struct {
	repo              repository.IUserRepository
	RefreshTokensRepo repository.IRefreshTokensRepo
}

func NewUserService(repo repository.IUserRepository, RefreshTokensRepo repository.IRefreshTokensRepo) UserService {
	return &userService{
		repo:              repo,
		RefreshTokensRepo: RefreshTokensRepo,
	}
}

// GetUserByID ...
func (u *userService) GetUserByID(id int) (*appresponse.UserResponse, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return &appresponse.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// CreateUser ...
func (u *userService) CreateUser(user *apprequest.UserRequest) (*appresponse.AuthenticateUserResponse, error) {
	const funcName = "service.CreateUser"
	in := userApprequestToDomain(*user)

	id, err := u.repo.CreateUser(in)
	if err != nil {
		log.Printf("%s: failed to create user: %v", funcName, err)
		return nil, err
	}

	// 3. Generate both access and refresh tokens (if using refresh tokens)
	accessToken, RefreshTokens, accessExpiry, refreshExpiry, err := helper.GenerateTokens(int(id), user.Email)
	if err != nil {
		return nil, err
	}

	// 4. Store the refresh token (recommended for modern flows)
	rt := &domain.RefreshTokens{
		UserID:    int(id),
		Token:     RefreshTokens,
		ExpiresAt: refreshExpiry,
		Revoked:   false,
	}
	if err := u.RefreshTokensRepo.SaveRefreshTokens(rt); err != nil {
		return nil, err
	}

	// 5. Return the response with tokens and expiries
	return &appresponse.AuthenticateUserResponse{
		Success:        true,
		AccessToken:    accessToken,
		RefreshTokens:  RefreshTokens,
		AccessExpires:  accessExpiry.Format(time.RFC3339),
		RefreshExpires: refreshExpiry.Format(time.RFC3339),
	}, nil
}

// userApprequestToDomain ...
func userApprequestToDomain(user apprequest.UserRequest) *domain.User {
	return &domain.User{
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
		Role:     user.Role, // Assuming Role is part of UserRequest
	}
}

// AuthenticateUser implements UserService.
func (u *userService) AuthenticateUser(auth *apprequest.AuthenticateUserRequest) (*appresponse.AuthenticateUserResponse, error) {
	user, err := u.repo.GetUserByEmail(auth.Email)
	if err != nil {
		return nil, err
	}
	if !helper.CheckPasswordHash(auth.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, RefreshTokens, accessExp, refreshExp, err := helper.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Save refresh token in DB
	rt := &domain.RefreshTokens{
		UserID:    user.ID,
		Token:     RefreshTokens,
		ExpiresAt: refreshExp,
		Revoked:   false,
	}
	if err := u.RefreshTokensRepo.SaveRefreshTokens(rt); err != nil {
		return nil, err
	}

	return &appresponse.AuthenticateUserResponse{
		Success:        true,
		AccessToken:    accessToken,
		AccessExpires:  accessExp.Format(time.RFC3339),
		RefreshTokens:  RefreshTokens,
		RefreshExpires: refreshExp.Format(time.RFC3339),
	}, nil
}

// ChangePassword ...
func (u *userService) ChangePassword(in *apprequest.ChangePasswordRequest) error {
	user, err := u.repo.GetUserByEmail(in.Email)
	if err != nil {
		return err
	}

	if !helper.CheckPasswordHash(in.CurrentPassword, user.Password) {
		return errors.New("invalid credentials")
	}

	err = u.repo.ChangePassword(user.ID, in.NewPassword)
	if err != nil {
		return err
	}

	return nil
}

// RefreshTokenss ...
func (u *userService) RefreshTokenss(oldRefreshTokens string) (*appresponse.AuthenticateUserResponse, error) {
	// 1. Get token from DB
	storedToken, err := u.RefreshTokensRepo.GetRefreshTokens(oldRefreshTokens)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. Check if it's revoked or expired
	if storedToken.Revoked || time.Now().After(storedToken.ExpiresAt) {
		return nil, errors.New("refresh token expired or revoked")
	}

	// 3. Get the user
	user, err := u.repo.GetUserByID(storedToken.UserID)
	if err != nil {
		return nil, err
	}

	// 4. (Optional) revoke the old refresh token to prevent reuse (token rotation)
	if err := u.RefreshTokensRepo.RevokeRefreshTokens(oldRefreshTokens); err != nil {
		return nil, err
	}

	// 5. Generate new tokens
	accessToken, newRefreshTokens, accessExp, refreshExp, err := helper.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// 6. Store the new refresh token
	rt := &domain.RefreshTokens{
		UserID:    user.ID,
		Token:     newRefreshTokens,
		ExpiresAt: refreshExp,
		Revoked:   false,
	}
	if err := u.RefreshTokensRepo.SaveRefreshTokens(rt); err != nil {
		return nil, err
	}

	// 7. Return the new tokens
	return &appresponse.AuthenticateUserResponse{
		Success:        true,
		AccessToken:    accessToken,
		AccessExpires:  accessExp.Format(time.RFC3339),
		RefreshTokens:  newRefreshTokens,
		RefreshExpires: refreshExp.Format(time.RFC3339),
	}, nil
}

// Logout ...
func (u *userService) Logout(RefreshTokens string) error {
	// Find token in DB and revoke it
	storedToken, err := u.RefreshTokensRepo.GetRefreshTokens(RefreshTokens)
	if err != nil {
		return errors.New("invalid refresh token")
	}
	if storedToken == nil {
		return errors.New("refresh token not found")
	}
	if storedToken.Revoked {
		return errors.New("refresh token already revoked")
	}
	// Revoke the token
	return u.RefreshTokensRepo.RevokeRefreshTokens(RefreshTokens)
}

// UpdateUserByID updates user details
func (u *userService) UpdateUserByID(id int, user *apprequest.UpsertUserPersonalDetailRequest) (*appresponse.UserResponse, error) {
	// Fetch existing user
	existingUser, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	// Update fields
	if user.Phone != "" {
		existingUser.Phone = user.Phone
	}
	if user.Role != "" {
		existingUser.Role = user.Role
	}
	if user.Avatar != "" {
		existingUser.Avatar = user.Avatar
	}
	// Save updated user
	if err := u.repo.UpdateUserByID(id, existingUser); err != nil {
		return nil, err
	}
	// Convert to response format
	response := &appresponse.UserResponse{
		ID:        existingUser.ID,
		Name:      existingUser.Name,
		Email:     existingUser.Email,
		Phone:     existingUser.Phone,
		Role:      existingUser.Role,
		Avatar:    existingUser.Avatar,
		IsActive:  existingUser.IsActive,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt,
	}
	return response, nil
}
