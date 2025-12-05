package shards

import (
	"context"
	"database/sql"
	"go-sharding-basic/internal/models"
	"log"
)

type Shard struct {
	db   *sql.DB
	name string
}

func NewShard(name string, dsn string) (*Shard, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Shard{
		db:   db,
		name: name,
	}, nil
}

func (s *Shard) SaveUser(ctx context.Context, username string, password string) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, username, password)
	if err != nil {
		return err
	}

	log.Println("Insert en shard:", s.name)

	return nil
}

func (s *Shard) GetUser(ctx context.Context, username string) (*models.User, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, username, password FROM users WHERE username = $1",
		username,
	)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}
