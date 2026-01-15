package httpclient

import (
	"context"
	"fmt"
	"micro-warehouse/merchant-service/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CachedUserClient struct {
	client UserClientInterface
	redis  *redis.RedisClient
	ttl    time.Duration
}

func NewCachedUserClient(userClient UserClientInterface, redisClient *redis.RedisClient) *CachedUserClient {
	return &CachedUserClient{
		client: userClient,
		redis:  redisClient,
		ttl:    1 * time.Hour,
	}
}

func (cpc *CachedUserClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("user:%s:%d", prefix, id)
}

func (cpc *CachedUserClient) generateCacheKeyMultiple(prefix string, ids []uint) string {
	key := fmt.Sprintf("user:%s", prefix)
	for _, id := range ids {
		key += fmt.Sprintf(":%d,", id)
	}
	return key[:len(key)-1]
}

func (cpc *CachedUserClient) GetUserByID(ctx context.Context, id uint) (*UserResponse, error) {
	cacheKey := cpc.generateCacheKey("single", id)

	var cachedUser UserResponse
	if err := cpc.redis.Get(ctx, cacheKey, &cachedUser); err != nil {
		log.Infof("[CachedUserClient] GetUserByID - 1: %v", err)
		return &cachedUser, nil
	}

	user, err := cpc.client.GetUserByID(ctx, id)
	if err != nil {
		log.Errorf("[CachedUserClient] GetUserByID - 2: %v", err)
		return nil, err
	}

	err = cpc.redis.Set(ctx, cacheKey, user, cpc.ttl)
	if err != nil {
		log.Errorf("[CachedUserClient] GetUserByID - 3: %v", err)
		return nil, err
	}

	return user, nil
}
