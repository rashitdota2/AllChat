package infrastructure

import "errors"

const (
	BadRequest  = "Неправильный запрос"
	ServerError = "InternalServerError"
)

var (
	ErrAlreadyExist  = errors.New("already exist")
	ErrIncorrectInfo = errors.New("not correct password or login")
)
