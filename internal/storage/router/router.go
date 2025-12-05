package router

import (
	"context"
	"crypto/sha256"
	"go-sharding-basic/internal/models"
)

type UserStorage interface {
	SaveUser(ctx context.Context, username string, password string) error
	GetUser(ctx context.Context, username string) (*models.User, error)
}

type ShardRouter struct {
	shards []UserStorage
}

func NewShardRouter(shards []UserStorage) *ShardRouter {
	return &ShardRouter{
		shards: shards,
	}
}

func (r *ShardRouter) SaveUser(ctx context.Context, username string, password string) error {
	shard := r.pickShard(username)
	return shard.SaveUser(ctx, username, password)
}

func (r *ShardRouter) GetUser(ctx context.Context, username string) (*models.User, error) {
	shard := r.pickShard(username)
	return shard.GetUser(ctx, username)
}

func (r *ShardRouter) pickShard(username string) UserStorage {
	h := sha256.Sum256([]byte(username))
	index := int(h[0]) % len(r.shards)
	return r.shards[index]
}
