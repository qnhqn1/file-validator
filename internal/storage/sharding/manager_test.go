package sharding

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stathat/consistent"
	"github.com/stretchr/testify/require"
)

func TestShardForKey(t *testing.T) {
	manager := &ShardManager{
		shards: make(map[string]*pgxpool.Pool),
		ring:   consistent.New(),
	}


	shardNames := []string{"shard0", "shard1", "shard2"}
	for _, name := range shardNames {
		manager.shards[name] = &pgxpool.Pool{} // dummy
		manager.ring.Add(name)
	}

	tests := []struct {
		key string
	}{
		{key: "a"},
		{key: "b"},
		{key: "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			pool := manager.ShardForKey(tt.key)
			require.NotNil(t, pool)

			pool2 := manager.ShardForKey(tt.key)
			require.Equal(t, pool, pool2)
		})
	}
}


