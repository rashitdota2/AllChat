package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
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
	if err = srv.rdc.Set(ctx, refresh, fmt.Sprintf("%d_%s", id, name), time.Hour*24*60).Err(); err != nil {
		srv.logger.Error("failed to set tokens", zap.Error(err))
		return "", "", err
	}
	return access, refresh, nil
}

func (srv *Service) Refresh(ctx context.Context, token string) (string, string, error) {
	info, err := srv.rdc.Get(ctx, token).Result()
	if err != nil {
		srv.logger.Error("failed to get token", zap.Error(err))
		return "", "", err
	}
	if info == "" {
		return "", "", jwt.ErrInvalidKey
	}
	user := strings.Split(info, "_")
	id, _ := strconv.Atoi(user[0])
	acs, rt, err := GetTokens(id, user[1])
	if err != nil {
		srv.logger.Error("failed to get tokens", zap.Error(err))
		return "", "", err
	}
	return acs, rt, nil
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

func (srv *Service) UpdAvatar(ctx context.Context, claims models.TokenClaims, img []byte) error {
	// need to add go
	paht := fmt.Sprintf("avatars/user_%d.png", claims.UserId)
	exist, err := srv.rdc.HExists(ctx, "avatars", fmt.Sprintf("%d", claims.UserId)).Result()
	if err != nil {
		srv.logger.Error("failed to check avatar on exist", zap.Error(err))
		return err
	}
	if !exist {
		if err = srv.rdc.HSet(ctx, "avatars", claims.UserId, paht).Err(); err != nil {
			srv.logger.Error("failed to set into rdc avatar", zap.Error(err))
			return err
		}
		file, err := os.Create(paht)
		defer file.Close()
		_, err = file.Write(img)
		if err != nil {
			srv.logger.Error("failed to write image to file", zap.Error(err))
			return err
		}
		err = srv.r.SetAvatar(ctx, claims.UserId, paht)
		if err != nil {
			srv.logger.Error("failed to set into psql avatar", zap.Error(err))
			return err
		}
		return nil
	}
	file, err := os.OpenFile(paht, os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		srv.logger.Error("failed to open file", zap.Error(err))
		return err
	}
	defer file.Close()
	if err = file.Truncate(0); err != nil {
		srv.logger.Error("failed to truncate file", zap.Error(err))
		return err
	}
	if _, err = file.Write(img); err != nil {
		srv.logger.Error("failed to write image to file", zap.Error(err))
		return err
	}
	return nil
}
