package api

import (
	"fmt"
	"time"

	"github.com/boscod/responsewatch-cli/internal/models"
)

// AuthAPI handles authentication-related API calls
type AuthAPI struct {
	Client *Client
}

// NewAuthAPI creates a new auth API
func NewAuthAPI(client *Client) *AuthAPI {
	return &AuthAPI{Client: client}
}

// Login authenticates a user and returns tokens
func (a *AuthAPI) Login(email, password string) (*models.LoginResponse, error) {
	req := models.LoginRequest{
		Email:    email,
		Password: password,
	}

	var resp models.LoginResponse
	if err := a.Client.Post("/auth/login", req, &resp, false); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Logout invalidates the current token
func (a *AuthAPI) Logout() error {
	return a.Client.Post("/auth/logout", nil, nil, true)
}

// Me gets the current user profile
func (a *AuthAPI) Me() (*models.User, error) {
	var user models.User
	if err := a.Client.Get("/auth/me", &user, true); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates the user profile
func (a *AuthAPI) UpdateProfile(req models.UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := a.Client.Put("/auth/profile", req, &user, true); err != nil {
		return nil, err
	}
	return &user, nil
}

// ChangePassword changes the user password
func (a *AuthAPI) ChangePassword(currentPassword, newPassword string) error {
	req := models.ChangePasswordRequest{
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	}
	return a.Client.Put("/auth/change-password", req, nil, true)
}

// SaveLogin saves login credentials to config
func (a *AuthAPI) SaveLogin(resp *models.LoginResponse) error {
	a.Client.Config.Auth.Token = resp.AccessToken
	a.Client.Config.Auth.RefreshToken = resp.RefreshToken
	a.Client.Config.Auth.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	a.Client.Config.User.Email = resp.User.Email
	a.Client.Config.User.Name = resp.User.FullName
	return a.Client.Config.Save()
}

// ClearAuth clears authentication data
func (a *AuthAPI) ClearAuth() error {
	a.Client.Config.ClearAuth()
	return a.Client.Config.Save()
}

// CheckAuth checks if the user is authenticated
func (a *AuthAPI) CheckAuth() error {
	if a.Client.Config.Auth.Token == "" {
		return fmt.Errorf("not authenticated. Please run 'rwcli login'")
	}

	if a.Client.Config.Auth.ExpiresAt.Before(time.Now()) {
		if a.Client.Config.Auth.RefreshToken == "" {
			return fmt.Errorf("session expired. Please run 'rwcli login'")
		}
		// Try to refresh
		if err := a.refreshToken(); err != nil {
			return fmt.Errorf("session expired. Please run 'rwcli login'")
		}
	}

	return nil
}

// refreshToken refreshes the access token
func (a *AuthAPI) refreshToken() error {
	url := "/auth/refresh"

	req := models.RefreshTokenRequest{
		RefreshToken: a.Client.Config.Auth.RefreshToken,
	}

	var resp models.RefreshTokenResponse
	if err := a.Client.Post(url, req, &resp, false); err != nil {
		return err
	}

	a.Client.Config.Auth.Token = resp.AccessToken
	a.Client.Config.Auth.RefreshToken = resp.RefreshToken
	a.Client.Config.Auth.ExpiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)

	return a.Client.Config.Save()
}
