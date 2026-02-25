# Kế hoạch & Chú thích chức năng toàn bộ File/Folder BE

Tài liệu này giải thích chi tiết mục đích, chức năng của từng thư mục và file code trong dự án `passiontech-beauty-contest` backend.

---

## 1. Top-Level (Thư mục gốc)

| File/Folder | Tác dụng |
|---|---|
| `cmd/` | Chứa điểm bắt đầu (entry point) của các chương trình (như server chính, CLI tool). |
| `configs/` | Chứa file mẫu cấu hình YAML (`passiontech-beauty-contest.yml`). Đây là nơi định nghĩa biến môi trường cho ứng dụng. |
| `deployments/` | Chứa các file để chạy hạ tầng (Docker, Kubernetes). Ví dụ: `docker-compose.yml` định nghĩa MongoDB, Redis ảo để dev. |
| `docs/` | Chứa tài liệu thiết kế dự án (Markdown), chú thích, quy hoạch thay đổi code. |
| `graph/` | Thư mục của GraphQL (`gqlgen`), chứa schema, interface và query/mutation resolver. |
| `internal/` | Chứa toàn bộ logic nghiệp vụ (domain logic) của ứng dụng. Go cấm thư mục bên ngoài import code từ `internal/`. |
| `scripts/` | Chứa các file bash/powershell chạy cài đặt, migration, dọn dẹp dự án. |

---

## 2. Thư mục `cmd/passiontech-beauty-contest`

Nơi khởi động ứng dụng thực tế.

| File | Tác dụng |
|---|---|
| `main.go` | File chạy đầu tiên (`func main()`). Chỉ có 1 việc: load biến cài đặt và gọi bước Khởi tạo. |
| `initial/init.go` | File "Khởi động thực tế": Kết nối Database, cấu hình Cache, khởi tạo Service, gọi Router và bắt đầu lắng nghe cổng 8080. |

---

## 3. Thư mục `internal/` (Core Logic)

Đây là nơi chứa 90% logic code dự án tuân theo chuẩn Phân tầng (Layered Architecture).

### `internal/model/` (Cấu trúc dữ liệu Database)
Chứa các Struct đại diện cho cấu trúc bảng/tài liệu trong CSDL MongoDB.
- `contestant.go`: Chứa model thí sinh (tuổi, chiều cao, sở thích...)
- `user.go`: Chứa model tài khoản người dùng (username, password hash)
- `score.go`: Chứa model điểm khảo thí (điểm các phần thi của BGK)
- `feedback.go` / `schedule.go`: Cấu trúc góp ý, lịch trình sự kiện.
- `constants.go`: Lưu các hằng số dùng chung toàn project (vd format ngày `YYYY-MM-DD`).

### `internal/types/` (Giao diện chuẩn - Interfaces)
Định nghĩa sẵn các "hợp đồng" (Interface) để các tầng nói chuyện với nhau mà không dính chặt vào công nghệ cụ thể (Dependency Injection).
- `auth.go` / `contestant.go` / `score.go`...: Định nghĩa các hàm mà *Service* và *Repository* bắt buộc phải có (Vd: `Create()`, `GetByID()`).
- `cache.go`: Định nghĩa các hàm gọi Cache chuẩn (`Set()`, `Get()`).

### `internal/dao/` (Data Access Object - Tầng CSDL)
Nơi duy nhất giao tiếp trực tiếp với MongoDB, thực thi các câu lệnh thêm, sửa, xóa.
- `contestant.go` / `user.go` / `score.go`...: Viết code gọi Mongo Driver (find, insertOne, updateOne).
- `storage.go`: Code để lưu trữ file (lưu ảnh vào máy chủ thay vì DB).

### `internal/service/` (Business Logic - Tầng Nghiệp vụ)
Tim của dự án. Thực hiện tính toán, kiểm tra điều kiện kinh doanh trước khi gọi DAO.
- `auth.go`: Logic đăng ký, đăng nhập. (Có đúng password không? Trả về JWT).
- `contestant.go`: Logic tạo hồ sơ thí sinh (Có đủ 18 tuổi không? Lọc những thí sinh bị cấm...).
- `score.go`: Tính tổng điểm, kiểm tra giám khảo có quyền chấm không.
- `feedback.go` / `schedule.go`: Logic nghiệp vụ phụ trợ.

### `internal/handler/` (Tầng Controller - Tiếp nhận Request URL / REST)
Nhận HTTP/REST API request (vd khi Frontend gửi JSON lên), trả về JSON.
- `auth.go`: Chuyển JSON thành Struct, gọi qua `AuthService`, sau đó gửi JSON "200 OK" về FE.
- `upload.go`: Tiếp nhận file upload từ Client (Multipart form) và lưu hình ảnh.
- `validation/`: Kiểm tra dữ liệu tĩnh (Email đúng chuẩn không? Số CCCD có bị thiếu số không?).

### `internal/routers/` (Định tuyến API)
- `router.go`: Kết nối các URL (như `/api/auth/login`) trỏ vào đúng Handler. Đặt cấu hình CORS và cấu hình thư mục lưu ảnh.

### `internal/server/` (Cấu hình máy chủ & Bảo mật base)
- `bcrypt.go`: Cung cấp hàm mã hóa mật khẩu (băm hash).
- `jwt.go`: Hàm tạo và giải mã chuỗi token (JWT) để giữ kết nối đăng nhập của User.

### `internal/cache/` (Kết nối Caching)
- `redis_adapter.go`: Cài đặt thực tế interface `CacheService` để dùng Redis tăng tốc độ (Vd: Nhớ số người online, rate limit block spam).

### `internal/config/` & `internal/ecode/`
- `config/config.go`: Đọc file `configs/passiontech-beauty-contest.yml`, biến thành struct tên là `Config` để code lấy ra port, pass mongdb an toàn.
- `ecode/ecode.go`: Quy tụ lại toàn bộ mã lỗi (Error codes), ví dụ: `1001: Lỗi Token`, `1002: Lỗi thiếu quyền`.

---

## 4. Thư mục `graph/` (GraphQL)

Dự án này sử dụng GraphQL (Gqlgen) thay cho API REST truyền thống trên một số tính năng phức tạp. Trình xử lý chính giống như các *Handler*.

| File/Folder | Tác dụng |
|---|---|
| `schema/` | Chứa các file đuôi `.graphqls`. Nơi bạn định nghĩa Query/Mutation/Type theo ngôn ngữ GraphQL (như 1 bản thiết kế nhà). |
| `generated/` | CODE AUTO TẠO. Không bao giờ sửa file này. Gqlgen sẽ tự sinh code dựa trên bản thiết kế schema. |
| `model/models_gen.go` | CODE AUTO TẠO. Các struct Go tự tạo ra để hứng dữ liệu GraphQL. |
| `resolver/` | Nơi dev TỰ VIẾT logic thật sự cho các node GraphQL (như Controller gọi tới Service). Cứ có node mới trong schema là thêm code ở đây. |
| `resolver/mapper.go` | Mapping: Biến `model.Contestant` (từ thư mục `internal/model` cho DB) thành `gqlmodel.Contestant` (để trả về Json GraphQL). |
| `middleware/auth.go` | Cổng bảo vệ GraphQL/HTTP. Chặn người dùng nếu họ gửi request không có chứa Token JWT hợp lệ. Bóc tách ra UserID. |
| `middleware/ratelimit.go` | Giới hạn IP spam server. Ví dụ 1 phút 1 user chỉ được ping API 100 lần. |

---

## Luồng chạy Data (Flow Logic)

Khi Frontend gọi 1 lượt "Sửa Hồ Sơ Thí Sinh":

1. FE gửi GraphQL Query (chứa JWT Token) 
2. 👉 Bị chặn lại ở `graph/middleware/auth.go` để lấy UserID
3. 👉 Request đi vào `graph/resolver/contestant.resolvers.go` (`UpdateContestantProfile` mutation)
4. 👉 Resolver gọi hàm `ToDomainUpdateContestant` trong `mapper.go`
5. 👉 Chạy qua `internal/handler/validation/contestant.go` kiểm tra tính logic
6. 👉 Đi sâu vào `internal/service/contestant.go` kiểm tra xem người này có đúng chủ hồ sơ không.
7. 👉 `service` gọi xuống `internal/dao/contestant.go`
8. 👉 `dao` chạy lệnh `UpdateOne` thẳng xuống Database MongoDB.
9. 👉 Dữ liệu mới trả ngược dòng lên thành format GraphQL cho FE tải lại lưới dữ liệu.
