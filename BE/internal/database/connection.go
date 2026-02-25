/*
Package database chịu trách nhiệm thiết lập kết nối đến tầng lưu trữ vật lý (MongoDB).
Tác dụng:
- Hàm Connect được gọi 1 lần duy nhất ở lúc khởi động (InitAndRun).
- Duy trì Connection Pool (hồ chứa kết nối) để tái sử dụng thay vì cứ mỗi request lại tạo mới kết nối.
*/
package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect tạo Pool HTTP tới MongoDB và ping thử để đảm bảo Database đang sống trước khi Server nhận request.
func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	// Production hardening:
	clientOptions.SetMinPoolSize(10)
	clientOptions.SetMaxPoolSize(100)
	clientOptions.SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
