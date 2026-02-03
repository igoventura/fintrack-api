package service

import (
	"context"

	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*dto.AuthResponse, error)
	Login(ctx context.Context, email, password string) (*dto.AuthResponse, error)
}

type SupabaseAuthService struct {
	client gotrue.Client
}

func NewSupabaseAuthService(projectID, apiKey string) *SupabaseAuthService {
	return &SupabaseAuthService{
		client: gotrue.New(projectID, apiKey),
	}
}

func (s *SupabaseAuthService) Register(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	req := types.SignupRequest{
		Email:    email,
		Password: password,
	}
	resp, err := s.client.Signup(req)
	if err != nil {
		return nil, err
	}

	// For signup, we might not get a session immediately if email confirmation is enabled.
	// But assuming we get a user.
	// If AccessToken is empty, it means confirmation is required.
	return &dto.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: dto.User{
			ID:    resp.User.ID.String(),
			Email: resp.User.Email,
		},
	}, nil
}

func (s *SupabaseAuthService) Login(ctx context.Context, email, password string) (*dto.AuthResponse, error) {
	resp, err := s.client.SignInWithEmailPassword(email, password)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		User: dto.User{
			ID:    resp.User.ID.String(),
			Email: resp.User.Email,
		},
	}, nil
}
