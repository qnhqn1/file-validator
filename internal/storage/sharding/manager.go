package sharding

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stathat/consistent"

	"github.com/qnhqn1/file-validator/config"
)


type ShardManager struct {
	shards map[string]*pgxpool.Pool
	ring   *consistent.Consistent
}


func NewShardManager(ctx context.Context, cfg config.DatabaseConfig) (*ShardManager, error) {
	if len(cfg.Shards) == 0 {
		return nil, fmt.Errorf("шарды не настроены")
	}
	m := &ShardManager{
		shards: make(map[string]*pgxpool.Pool),
		ring:   consistent.New(),
	}
	for _, s := range cfg.Shards {
		pool, err := newShardPool(ctx, s.ConnURL)
		if err != nil {
			m.Close()
			return nil, fmt.Errorf("создать шард %s: %w", s.Name, err)
		}
		m.shards[s.Name] = pool
		m.ring.Add(s.Name)
	}
	return m, nil
}

func newShardPool(ctx context.Context, connURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connURL)
	if err != nil {
		return nil, fmt.Errorf("разобрать pg config: %w", err)
	}
	cfg.MaxConns = 10
	cfg.MinConns = 1
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("создать пул: %w", err)
	}
	return pool, nil
}


func (m *ShardManager) Close() {
	for _, p := range m.shards {
		if p != nil {
			p.Close()
		}
	}
}


func (m *ShardManager) ShardForKey(key string) *pgxpool.Pool {
	shardName, err := m.ring.Get(key)
	if err != nil {

		return m.Primary()
	}
	return m.shards[shardName]
}


func (m *ShardManager) Primary() *pgxpool.Pool {

	for _, pool := range m.shards {
		return pool // возвращаем любой, но лучше первый
	}
	return nil
}


