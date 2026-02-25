/*
Package main là điểm khởi đầu (entry point) của toàn bộ ứng dụng Backend.
Theo chuẩn của dự án passiontech-sample, thư mục cmd/ chỉ chứa file main cực kỳ mỏng.
Tác dụng duy nhất: Setup môi trường (.env) và nhường quyền cho thư mục initial/ để khởi động.
Tuyệt đối KHÔNG viết logic nghiệp vụ (database, xử lý request) ở đây.
*/
package main

import (
	"cuoc_thi_hoa_hau/cmd/passiontech-beauty-contest/initial"

	"github.com/joho/godotenv"
)

func main() {
	// Bước 1: Load file .env (nếu có) để nạp các biến môi trường vào hệ thống
	_ = godotenv.Load() // Load env vars if .env exists

	// Bước 2: Gọi hàm InitAndRun để bắt đầu quá trình Dependency Injection và bật Web Server
	initial.InitAndRun()
}
