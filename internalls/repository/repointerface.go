package repository

import (
	"context"
	"workwithimages/domain/models"
)

type RepoInterface interface {
	Registration(ctx context.Context, info models.SignIn) error
	CheckLogin(ctx context.Context, auth models.Auth) (*models.UserProfile, string, error)
	GetProfile(ctx context.Context, userId int) (models.UserProfile, error)
	UpdateProfile(ctx context.Context, profile models.UserProfile, userId int) error
	SetAvatar(ctx context.Context, id int, path string) error
}
