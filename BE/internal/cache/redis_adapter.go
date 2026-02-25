/*
Package cache là Tầng Hạ Tầng (Infrastructure implementation) cho bộ nhớ đệm.
Tác dụng:
- Cài đặt Interface CacheService (trong package types) bằng công nghệ Redis.
- Giúp ứng dụng ghi/đọc RAM thay vì SSD (DB), giúp truy vấn cực nhanh (Microsecond).
- Dùng để: Đếm Rate Limit (chống DDoS), lưu danh sách token đăng xuất, v.v.
*/
package cache

import (
	"context"
	"encoding/json"
	"time"

	"cuoc_thi_hoa_hau/internal/types"

	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client *redis.Client
}

// NewRedisAdapter khởi tạo một con trỏ tới CSDL Redis.
func NewRedisAdapter(addr string, password string, db int) types.CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisAdapter{
		client: rdb,
	}
}

// Set lấy cái Struct của Go, biến thành chuỗi JSON (Marshal) rồi nhét vào bộ nhớ phân mảnh Redis.
func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonVal, expiration).Err()
}

func (r *RedisAdapter) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisAdapter) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}
