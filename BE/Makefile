.PHONY: run build build-docker clean generate

# Lệnh chạy Backend trực tiếp trên máy
run:
	go run ./cmd/passiontech-beauty-contest/

# Lệnh biên dịch ra file nhị phân
build:
	go build -o bin/server ./cmd/passiontech-beauty-contest/

# Lệnh sinh code GraphQL tự động (gqlgen)
generate:
	go run github.com/99designs/gqlgen generate

# Dọn dẹp thư mục build
clean:
	rm -rf bin/

# Build Docker Image nội bộ cho BE
build-docker:
	docker build -t beauty-contest-be .
