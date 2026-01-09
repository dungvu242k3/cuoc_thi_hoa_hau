# Hướng dẫn chạy dự án Cuộc thi Hoa hậu (Beauty Contest)

Tài liệu này hướng dẫn cách cài đặt và chạy dự án trên máy cục bộ.

## 1. Yêu cầu hệ thống
Đảm bảo máy tính của bạn đã cài đặt:
- **Go** (Golang): Để chạy Backend. [Tải về tại đây](https://go.dev/dl/).
- **Node.js** và **npm**: Để chạy Frontend. [Tải về tại đây](https://nodejs.org/).

## 2. Cấu trúc dự án
- `BE/`: Thư mục chứa mã nguồn Backend (viết bằng Go).
- `FE/`: Thư mục chứa mã nguồn Frontend (viết bằng React/TypeScript).

## 3. Cách chạy dự án

Bạn cần mở **2 cửa sổ Terminal** (Command Prompt hoặc PowerShell) riêng biệt để chạy song song Backend và Frontend.

### Phần 1: Chạy Backend (API)
Backend sẽ xử lý dữ liệu và cung cấp API cho Frontend.

1.  Mở Terminal thứ nhất.
2.  Di chuyển vào thư mục `BE`:
    ```bash
    cd BE
    ```
3.  Chạy lệnh khởi động server:
    ```bash
    go run cmd/api/main.go
    ```
    *(Nếu hệ thống báo thiếu thư viện, hãy chạy `go mod tidy` trước khi chạy lệnh trên)*.

=> **Thành công:** Backend sẽ hoạt động (thường mặc định tại cổng 8080).

### Phần 2: Chạy Frontend (Giao diện Web)
Frontend là giao diện trang web mà bạn sẽ tương tác.

1.  Mở Terminal thứ hai.
2.  Di chuyển vào thư mục `FE`:
    ```bash
    cd FE
    ```
3.  Cài đặt các thư viện cần thiết (chỉ cần làm lần đầu):
    ```bash
    npm install
    ```
4.  Chạy lệnh khởi động trang web:
    ```bash
    npm run dev
    ```

=> **Thành công:** Terminal sẽ hiện đường dẫn truy cập, thường là `http://localhost:5173`. Bạn hãy mở đường dẫn này trên trình duyệt.

## 4. Tóm tắt lệnh

| Thành phần | Thư mục | Lệnh chạy | Địa chỉ mặc định |
| :--- | :--- | :--- | :--- |
| **Backend** | `BE` | `go run cmd/api/main.go` | `http://localhost:8080` |
| **Frontend** | `FE` | `npm run dev` | `http://localhost:5173` |
