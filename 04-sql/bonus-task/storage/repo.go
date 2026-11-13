package storage

import (
	"context"
	"errors"

	"github.com/verdex4/04-sql-bonus/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Storage — хранилище items
type Storage interface {
	Get(ctx context.Context, id int64) (domain.Item, error)
	Add(ctx context.Context, item domain.Item) (int64, error)
}
type Postgres struct {
	pool *pgxpool.Pool
}

// New создаёт пул соединений и проверяет, что база жива.
func New(ctx context.Context, dsn string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &Postgres{pool: pool}, nil
}

// Close освобождает соединения. Вызывается при остановке приложения.
func (s *Postgres) Close() {
	s.pool.Close()
}

const getByID = `SELECT title, description FROM items WHERE id = $1`

func (s *Postgres) Get(ctx context.Context, id int64) (domain.Item, error) {
	var it domain.Item
	err := s.pool.QueryRow(ctx, getByID, id).Scan(&it.Title, &it.Description)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.Item{}, status.Error(codes.NotFound, "item not found")
	case err != nil:
		return domain.Item{}, status.Error(codes.Internal, "db error")
	}
	return it, nil
}

const insert = `INSERT INTO items (title, description) VALUES ($1, $2) RETURNING id`

func (s *Postgres) Add(ctx context.Context, item domain.Item) (int64, error) {
	var id int64
	if err := s.pool.QueryRow(ctx, insert, item.Title, item.Description).Scan(&id); err != nil {
		return 0, status.Error(codes.Internal, "db error")
	}
	return id, nil
}

func (s *Postgres) OldestInCity(ctx context.Context, city string) (domain.User, error) {
	q := "SELECT name, age FROM users WHERE age = (SELECT MAX(age) as m FROM users);"
	var u domain.User
	err := s.pool.QueryRow(ctx, q).Scan(&u.Name, &u.Age)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return domain.User{}, status.Error(codes.NotFound, "oldest users not found")
	case err != nil:
		return domain.User{}, status.Error(codes.Internal, "db error")
	}
	return u, nil
}
