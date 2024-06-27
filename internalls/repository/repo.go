package repository

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"workwithimages/domain/models"
)

type Repository struct {
	db  *pgxpool.Pool
	rdc *redis.Client
}

func NewRepo(db *pgxpool.Pool, rdc *redis.Client) RepoInterface {
	return &Repository{
		db:  db,
		rdc: rdc,
	}
}

func (r *Repository) Registration(ctx context.Context, info models.SignIn) error {

	//If You want to use redis to store acc info instead of psql
	//(OPTIONALLY: Use tx after exist part)
	//
	//exists, err := r.rdc.HExists(ctx, "users_table", info.Login).Result()
	//if err != nil {
	//
	//	return err
	//}
	//if exists {
	//	return infrastructure.ErrAlreadyExist
	//}
	//id, err := r.rdc.Incr(ctx, "users").Result()
	//idstr := fmt.Sprintf("%d", id)
	//if err != nil {
	//	return err
	//}
	//err = r.rdc.HSet(ctx, "users_table", info.Login, idstr, "name_"+idstr, info.Name, "key_user_"+idstr, info.Key).Err()
	//if err != nil {
	//	return err
	//}
	if _, err := r.db.Exec(ctx, "insert into users (login,key,name) values ($1,$2,$3)", info.Login, info.Key, info.Name); err != nil {
		return err
	}
	return nil
}

func (r *Repository) CheckLogin(ctx context.Context, info models.Auth) (*models.UserProfile, string, error) {
	var user models.UserProfile
	var key string
	if err := r.db.QueryRow(ctx, "select id, name, key from users where login = $1", info.Login).Scan(&user.Id, &user.Name, &key); err != nil {
		return nil, "", err
	}
	return &user, key, nil
}

func (r *Repository) GetProfile(ctx context.Context, userId int) (models.UserProfile, error) {
	var user models.UserProfile

	if err := r.db.QueryRow(ctx, "select id, name, description, avatar from users where id = $1", userId).Scan(&user.Id, &user.Name, &user.Description, &user.Avatar); err != nil {
		return models.UserProfile{}, err
	}
	return user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, profile models.UserProfile, userId int) error {
	if _, err := r.db.Exec(ctx, "update users set name = $1, description = $2 where id = $3", profile.Name, profile.Description, userId); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SetAvatar(ctx context.Context, id int, path string) error {
	if _, err := r.db.Exec(ctx, "update users set avatar = $1 where id = $2", path, id); err != nil {
		return err
	}
	return nil
}
