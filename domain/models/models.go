package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	AccessTime  = time.Hour * 1
	RefreshTime = time.Hour * 24 * 90
	RedisTokens = "tokens"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
}

type Auth struct {
	Login string `json:"login"`
	Key   string `json:"key"`
}

type SignIn struct {
	Name  string `json:"name"`
	Login string `json:"login"`
	Key   string `json:"key"`
}

type UserProfile struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
