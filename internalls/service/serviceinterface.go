package service

import (
	"context"
	"workwithimages/domain/models"
)

type ServiceInterface interface {
	AuthService
	GetProfile(ctx context.Context, userId int) (models.UserProfile, error)
	UpdateProfile(ctx context.Context, profile models.UserProfile, userId int) error
	UpdAvatar(ctx context.Context, claims models.TokenClaims, img []byte, ext string) error
}

type AuthService interface {
	Sign(ctx context.Context, info models.SignIn) error
	Login(ctx context.Context, auth models.Auth) (string, string, error)
	Refresh(ctx context.Context, token string) (string, string, error)
}
