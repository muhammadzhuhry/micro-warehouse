package httpclient

import (
	"context"
	"fmt"
	"micro-warehouse/merchant-service/pkg/redis"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type CachedWarehouseClient struct {
	client WarehouseClientInterface
	redis  *redis.RedisClient
	ttl    time.Duration
}

func NewCachedWarehouseClient(warehouseClient WarehouseClientInterface, redisClient *redis.RedisClient) *CachedWarehouseClient {
	return &CachedWarehouseClient{
		client: warehouseClient,
		redis:  redisClient,
		ttl:    1 * time.Hour,
	}
}

func (cpc *CachedWarehouseClient) generateCacheKey(prefix string, id uint) string {
	return fmt.Sprintf("warehouse:%s:%d", prefix, id)
}

func (cpc *CachedWarehouseClient) GetWarehouseByID(ctx context.Context, warehouseID uint) (*WarehouseResponse, error) {
	cacheKey := cpc.generateCacheKey("single", warehouseID)

	var cachedWarehouse WarehouseResponse
	if err := cpc.redis.Get(ctx, cacheKey, &cachedWarehouse); err == nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseByID - 1: %v", err)
		return &cachedWarehouse, nil
	}

	warehouse, err := cpc.client.GetWarehouseByID(ctx, warehouseID)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseByID - 2: %v", err)
		return nil, err
	}

	err = cpc.redis.Set(ctx, cacheKey, warehouse, cpc.ttl)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseByID - 3: %v", err)
		return nil, err
	}

	return warehouse, nil
}

func (cpc *CachedWarehouseClient) GetWarehouseProductStock(ctx context.Context, warehouseID, productID uint) (*WarehouseProductStockResponse, error) {
	cacheKey := cpc.generateCacheKey("single", warehouseID)

	var cachedStock WarehouseProductStockResponse
	if err := cpc.redis.Get(ctx, cacheKey, &cachedStock); err == nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseProductStock - 1: %v", err)
		return &cachedStock, nil
	}

	stock, err := cpc.client.GetWarehouseProductStock(ctx, warehouseID, productID)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseProductStock - 2: %v", err)
		return nil, err
	}

	err = cpc.redis.Set(ctx, cacheKey, stock, cpc.ttl)
	if err != nil {
		log.Errorf("[CachedWarehouseClient] GetWarehouseProductStock - 3: %v", err)
		return nil, err
	}

	return stock, nil
}
