# Hướng dẫn chạy dự án Cuộc thi Hoa hậu (Beauty Contest)

Dự án này là một hệ thống Fullstack bao gồm **Backend (Go + MongoDB + Redis)** và **Frontend (React/Vite)**. Có 2 cách chính để khởi chạy dự án: Chạy thủ công trên Terminal hoặc Chạy tất cả với Docker.

---

## Cách 1: Chạy thủ công trên máy (Chạy từng dịch vụ)

Cách này phù hợp nhất dành cho Developers muốn viết code và thử nghiệm thay đổi ngay lập tức trên máy cục bộ.

### Yêu cầu cài đặt
- **Go 1.24+**: Cho Backend.
- **Node.js (v20+) & npm**: Cho Frontend.
- **MongoDB**: Chạy cục bộ ở cổng `27017`.
- **Redis**: Chạy cục bộ ở cổng `6379`.

### Bước 1: Khởi động Backend
1. Mở Terminal mới và đi vào thư mục Backend:
   ```bash
   cd BE
   ```
2. Cấu hình biến môi trường:
   Đảm bảo bạn đã có file `BE/.env` với nội dung chuẩn (sao chép nội dung gốc nếu chưa có):
   ```ini
   MONGODB_URI=mongodb://localhost:27017
   JWT_SECRET=thay-doi-khoa-bao-mat-nay
   REDIS_ADDR=localhost:6379
   ```
3. Khởi động API Server:
   **Cách 1 (Lệnh chuẩn):**
   ```bash
   go run ./cmd/passiontech-beauty-contest/
   ```
   **Cách 2 (Dùng lệnh rút gọn qua Makefile - Khuyên dùng):**
   ```bash
   make run
   ```
   > **Lưu ý:** Nếu code báo lỗi thư viện, hãy chạy `go mod tidy` trước. Backend khởi động thành công sẽ hiển thị ở cổng `http://localhost:8080`.

### Phụ lục: Các lệnh rút gọn (Makefile) cho Backend
Nếu bạn đang ở trong thư mục `BE/`, thay vì gõ dòng lệnh dài, bạn có thể dùng các lệnh `make` tiện lợi sau:
- `make run`: Khởi động Backend lên.
- `make build`: Biên dịch mã nguồn ra file chạy được vào thư mục `bin/`.
- `make generate`: Sinh lại code tự động cho GraphQL (gqlgen) khi bạn vừa thay đổi schema.
- `make clean`: Xóa nhanh các file build cũ.


### Bước 2: Khởi động Frontend
1. Mở một Terminal khác và đi vào thư mục Frontend:
   ```bash
   cd FE
   ```
2. Cài đặt các gói thư viện (Chỉ làm lần đầu tiên hoặc khi có thay đổi package.json):
   ```bash
   npm install
   ```
3. Khởi động giao diện Web:
   ```bash
   npm run dev
   ```
   > Frontend khởi động thành công sẽ hiển thị ở `http://localhost:5173`. Bạn có thể truy cập vào đường dẫn này bằng trình duyệt.

---

## Cách 2: Chạy toàn bộ hệ thống bằng Docker Compose (Khuyên dùng)

Cách này vô cùng tiện lợi, dành cho những ai **chỉ muốn chạy lên xem thử** hoặc chạy trên môi trường Product/Staging mà không cần cài đặt Go, Node, Mongo hay Redis rườm rà.

### Yêu cầu cài đặt
- **Docker Desktop** (Đã bật và đang chạy nền).

### Các bước khởi chạy
1. Mở Terminal ngay tại thư mục **Gốc** của dự án (nơi chứa file `docker-compose.yml`).
2. Gõ lệnh xây dựng và khởi chạy tất cả:
   ```bash
   docker-compose up -d --build
   ```
3. Đợi vài phút để hệ thống tải thư viện và biên dịch. Khi nhìn thấy các services chuyển sang trạng thái `Running` là thành công.

### Thông số các cổng sau khi chạy Docker
- **Frontend (Web):** `http://localhost:3000`
- **Backend (API + GraphQL):** `http://localhost:8080/graphql`
- **MongoDB:** `localhost:27017`
- **Redis:** `localhost:6379`

### Tắt hệ thống Docker
Khi không muốn chạy nữa, bạn gõ lệnh sau để dừng và dọn dẹp bộ nhớ:
```bash
docker-compose down
```

---

## 🛠 Tóm tắt cấu trúc cổng (Ports)

| Dịch Vụ | Công Nghệ | Port Local Manual | Port Docker Compose |
| :--- | :--- | :--- | :--- |
| **Backend API** | Go 1.24 | `8080` | `8080` |
| **Frontend UI** | React/Vite | `5173` | `3000` |
| **Database** | MongoDB | `27017` | `27017` |
| **Cache Server** | Redis | `6379` | `6379` |
