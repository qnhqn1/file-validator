package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qnhqn1/file-validator/internal/storage/sharding"
)


type StorageInterface interface {
	InsertEvent(ctx context.Context, key string, payload []byte) error
	PrimaryPool() *pgxpool.Pool
	Close()
}


type Storage struct {
	manager *sharding.ShardManager
}


func New(ctx context.Context, manager *sharding.ShardManager) (*Storage, error) {
	if manager == nil {
		return nil, fmt.Errorf("менеджер шардов nil")
	}
	return &Storage{manager: manager}, nil
}


func (s *Storage) Close() {
	if s.manager != nil {
		s.manager.Close()
	}
}


func (s *Storage) InsertEvent(ctx context.Context, key string, payload []byte) error {
	pool := s.manager.ShardForKey(key)
	if pool == nil {
		return fmt.Errorf("нет шарда для ключа")
	}
	_, err := pool.Exec(ctx, "INSERT INTO validator_events (key, payload) VALUES ($1, $2)", key, payload)
	return err
}


func (s *Storage) PrimaryPool() *pgxpool.Pool { return s.manager.Primary() }


