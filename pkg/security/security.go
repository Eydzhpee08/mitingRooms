package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

)

//Service Authorization
type Service struct {
	pool *pgxpool.Pool
}

var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInternal = errors.New("internal error")
var ErrExpireToken = errors.New("token expired")

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// Авторизация в этом методе, мы проверяем логин и пароль, если верны, то возвращаем true, если нет то false
// Генерирует токен пользователя
//Если пользователь не найден возрашает ErrNoSuchUser


func (s *Service) TokenForUsers(ctx context.Context, phone string, password string) (token string, err error) {
	var hash string
	var id int64

	err = s.pool.QueryRow(ctx, `SELECT id,password FROM users WHERE phone =$1`, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}
	if err != nil {
		return "", ErrInternal
	}

//bcrypt- генерируется автоматически и хранится в строке вместе с хешом) и проверки пароля:
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)

	n, err := rand.Read(buffer)

	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO users_tokens(token,users_id) VALUES($1,$2)`, token, id)
	if err != nil {
		return "", ErrInternal
	}

	return token, nil
}

func (s *Service) AuthenticateUsers(ctx context.Context, token string) (int64, error) {
	var id int64
	var expire time.Time

	err := s.pool.QueryRow(ctx, `SELECT users_id, expire FROM users_tokens WHERE token =$1`, token).Scan(&id, &expire)
	if err == pgx.ErrNoRows {
		log.Print(err)
		return 0, ErrNoSuchUser

	}
	if err != nil {
		log.Print(err)
		return 0, ErrInternal
	}

	tN := time.Now().Unix()
	log.Print(tN)
	tE := expire.Unix()
	log.Print(tE)

	if tN > tE {
		return 0, ErrExpireToken
	}

	return id, nil
}
