package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"os"
	"workwithimages/domain/infrastructure"
	"workwithimages/domain/models"
	"workwithimages/internalls/repository"
)

type Service struct {
	r      repository.RepoInterface
	logger *zap.Logger
	rdc    *redis.Client
}

func NewService(r repository.RepoInterface, logger *zap.Logger, rdc *redis.Client) ServiceInterface {
	return &Service{
		r:      r,
		logger: logger,
		rdc:    rdc,
	}
}

func (srv *Service) Sign(ctx context.Context, info models.SignIn) error {
	if err := srv.r.Registration(ctx, info); err != nil {
		//need to add check err on pgx err uniq
		srv.logger.Error("failed to register user", zap.Error(err))
		return err
	}
	return nil
}

func (srv *Service) Login(ctx context.Context, info models.Auth) (string, string, error) {
	id, name, err := srv.r.CheckLogin(ctx, info)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return "", "", infrastructure.ErrIncorrectInfo
	case err != nil:
		srv.logger.Error("failed to check user", zap.Error(err))
		return "", "", err
	}
	access, refresh, err := GetTokens(id, name)
	if err != nil {
		srv.logger.Error("failed to gen tokens", zap.Error(err))
		return "", "", err
	}

	return access, refresh, nil
}

func (srv *Service) GetProfile(ctx context.Context, userId int) (models.UserProfile, error) {

	profile, err := srv.r.GetProfile(ctx, userId)
	if err != nil {
		srv.logger.Error("failed to get profile", zap.Error(err))
		return models.UserProfile{}, err
	}
	return profile, nil
}

func (srv *Service) UpdateProfile(ctx context.Context, profile models.UserProfile, userId int) error {
	if err := srv.r.UpdateProfile(ctx, profile, userId); err != nil {
		srv.logger.Error("failed to update profile", zap.Error(err))
		return err
	}
	return nil
}

func (srv *Service) UpdAvatar(ctx context.Context, claims models.TokenClaims, img []byte, ext string) error {
	//need convert image to another format
	// need to add go
	id := fmt.Sprintf("%d", claims.UserId)
	if exist, _ := srv.rdc.HExists(ctx, "avatars", id).Result(); exist {
		oldPath, err := srv.rdc.HGet(ctx, "avatars", id).Result()
		if err != nil {
			srv.logger.Error("failed to get avatar", zap.Error(err))
			return err
		}
		err = os.Remove(oldPath)
		if err != nil {
			srv.logger.Error("failed to remove old avatar", zap.Error(err))
			return err
		}
	}

	paht := fmt.Sprintf("avatars/user_%d.%s", claims.UserId, ext)

	file, err := os.OpenFile(paht, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		srv.logger.Error("failed to open file", zap.Error(err))
		return err
	}
	defer file.Close()
	if _, err = file.Write(img); err != nil {
		srv.logger.Error("failed to write file", zap.Error(err))
		return err
	}
	if err = srv.rdc.HSet(ctx, "avatars", claims.UserId, paht).Err(); err != nil {
		srv.logger.Error("failed to set avatar", zap.Error(err))
		return err
	}
	if err = srv.r.SetAvatar(ctx, claims.UserId, paht); err != nil {
		srv.logger.Error("failed to set avatar", zap.Error(err))
		return err
	}
	return nil
}
