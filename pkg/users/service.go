package users

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

)

var ErrNotFound = errors.New("item not found")

var ErrInternal = errors.New("internal error")

type Service struct {
	pool  *pgxpool.Pool
}
func NewService(pool  *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Users struct {
	ID      	int64     `json:"id"`
	Name    	string    `json:"name"`
	Phone   	string    `json:"phone"`
	Password 	string	  `json:"password"`
	Active		bool 	  `json:"active"`
	Created 	time.Time `json:"created"`
}

func (s *Service) ByID(ctx context.Context, id int64) (*Users, error) {
	item := &Users{}

	err := s.pool.QueryRow(ctx, `
		SELECT id, name, phone, active, created FROM users WHERE id = $1
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
func (s *Service) All(ctx context.Context) (items []*Users, err error) {

	rows, err:= s.pool.Query(ctx, `
		SELECT * FROM users
	`)

	for rows.Next(){
		item := &Users{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
		if err != nil {
			log.Print(err)
		}

		items = append(items, item)
	}
	return items, nil
}
func (s *Service) AllActive(ctx context.Context) (items []*Users, err  error) {

	rows, err:= s.pool.Query(ctx, `
		SELECT * FROM users WHERE active
	`)

	for rows.Next(){
		item := &Users{}
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)
		if err != nil {
			log.Print(err)
		}

		items = append(items, item)
	}
	return items, nil
}

// //Save method
func (s *Service) Save(ctx context.Context, users *Users) (c *Users, err error) {

	item := &Users{}

	if users.ID == 0 {
		sqlStatement := `insert into users(name, phone, password) values($1, $2, $3) returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, users.Name, users.Phone, users.Password).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password,
			&item.Active,
			&item.Created)
	} else {
		sqlStatement := `update users set name=$1, phone=$2, password=$3 where id=$4 returning *`
		err = s.pool.QueryRow(ctx, sqlStatement, users.Name, users.Phone, users.Password, users.ID).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Password,
			&item.Active,
			&item.Created)
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}
	return item, nil

}

func (s *Service) RemoveById(ctx context.Context, id int64) (*Users, error) {
	item := &Users{}
	err := s.pool.QueryRow(ctx, `
	DELETE FROM users WHERE id=$1 RETURNING id,name,phone,active,created 
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}

func (s *Service) BlockByID(ctx context.Context, id int64) (*Users, error) {
	item := &Users{}
	err := s.pool.QueryRow(ctx, `
		UPDATE users SET active = false WHERE id = $1 RETURNING id, name, phone, active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
func (s *Service) UnBlockByID(ctx context.Context, id int64) (*Users, error) {
	item := &Users{}
	err := s.pool.QueryRow(ctx, `
		UPDATE users SET active = true WHERE id = $1 RETURNING id, name, phone,active, created
	`, id).Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil

}
