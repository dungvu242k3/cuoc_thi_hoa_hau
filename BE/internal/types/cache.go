/*
Package types định nghĩa hợp đồng cho bộ nhớ đệm (Cache).
Tác dụng:
- Tách biệt logic caching khỏi thư viện cụ thể (như go-redis).
- Tầng Service gọi Set/Get/Delete không cần biết cấu hình kết nối Redis ở dưới.
*/
package types

import (
	"context"
	"time"
)

type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
}
