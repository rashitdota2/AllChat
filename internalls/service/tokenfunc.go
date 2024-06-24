package service

import (
	"github.com/dgrijalva/jwt-go"
	"sync"
	"time"
	"workwithimages/domain/models"
)

func GetTokens(id int, name string) (string, string, error) {
	var (
		access, refresh string
		err, temperr    error
		wg              sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		access, temperr = AccessToken(id, name)
		if temperr != nil {
			err = temperr
		}
	}()
	go func() {
		defer wg.Done()
		refresh, temperr = RefreshToken(id)
		if temperr != nil {
			err = temperr
		}
	}()
	wg.Wait()
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func AccessToken(id int, name string) (string, error) {
	claims := models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(models.AccessTime).Unix(),
		},
		UserId: id,
		Name:   name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	acs, err := token.SignedString([]byte("secret"))
	return acs, err
}

func RefreshToken(id int) (string, error) {
	claims := models.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(models.RefreshTime).Unix(),
		},
		UserId: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	return token.SignedString([]byte("secret"))
}
