# Kế hoạch Chú thích toàn bộ Code (Annotate Codebase)

## 1. Overview (Tổng quan)
- **Mục tiêu**: Chú thích (comments) giải thích luồng hoạt động, pattern và tác dụng của từng lớp kiến trúc (layer) trực tiếp vào toàn bộ các file mã nguồn `.go` còn lại trong dự án. Giúp lập trình viên đọc code như đang được một Senior Reviewer hướng dẫn.
- **Phạm vi**: Những file cốt lõi ở `cmd`, `config`, `routers`, `handler/auth.go`, `service/auth.go`, `dao/user.go`, `types/user.go` đã được chú thích. Kế hoạch này tập trung vào số file còn lại trong thư mục `internal/` và `graph/`.

## 2. Project Type
- **BACKEND** (Golang, MongoDB, GraphQL, REST, Redis)

## 3. Success Criteria (Tiêu chí thành công)
- Toàn bộ các gói (packages) và các hàm (functions) lõi đều có GoDoc comment tiếng Việt giải thích rõ logic.
- Không thay đổi hành vi (behavior) của ứng dụng, chỉ thêm comments.
- Mã nguồn giữ được tính sạch sẽ, comments không bị dài dòng (spam).

## 4. Tech Stack
- Ngôn ngữ: Go 1.24
- Graph: gqlgen
- Web Framework: go-chi
- Database: MongoDB

## 5. File Structure Focus
Sẽ tập trung bổ sung comments dàn trải theo các thư mục sau:
- `internal/model/` và `internal/types/` (Hợp đồng và Cấu trúc DB)
- `internal/dao/` (Tầng kết nối CSDL thao tác CRUD)
- `internal/service/` (Tầng nghiệp vụ lõi)
- `internal/handler/` và `internal/server/` (Tầng HTTP và Tiện ích)
- `graph/` (Tầng GraphQL resolvers và middleware)

## 6. Task Breakdown

| Task ID | Component/Folder | Description (INPUT → OUTPUT → VERIFY) | Agent | Skill |
|---|---|---|---|---|
| 1 | `internal/model/` | Thêm mô tả các tags `bson` và `json` kết nối với DB như thế nào. <br/> **Verify**: Files `contestant.go`, `score.go`, ... có chú thích. | `backend-specialist` | `clean-code` |
| 2 | `internal/types/` | Chú thích tác dụng của Dependency Inversion qua Interfaces. <br/> **Verify**: Files `contestant.go`, `score.go`, `security.go` có chú thích. | `backend-specialist` | `clean-code` |
| 3 | `internal/dao/` | Chú thích cách gọi mongo-driver (InsertOne, Upsert, GetByID) và context timeout. <br/> **Verify**: Các files repo còn lại đều có chú thích tác dụng truy vấn CSDL. | `database-architect` | `database-design` |
| 4 | `internal/service/` | Giải thích logic nghiệp vụ phức tạp (như tính điểm, upload ảnh, logic phân trang). <br/> **Verify**: `contestant.go`, `score.go`, `feedback.go`, `schedule.go` có chú thích lõi. | `backend-specialist` | `clean-code` |
| 5 | `internal/handler/` & `server/` | Giải thích file upload xử lý multipart/form-data, bcrypt và JWT claims. <br/> **Verify**: `upload.go`, `jwt.go`, `bcrypt.go` có giải thích đầy đủ. | `backend-specialist` | `clean-code` |
| 6 | `graph/` | Thêm comment vào `mapper.go`, `resolvers` và `middleware` để giải thích luồng bảo vệ và xử lý của GraphQL. <br/> **Verify**: Middleware chặn Auth và resolver gọi Service được chú thích rõ ràng. | `backend-specialist` | `clean-code` |

## 7. Phase X: Verification
- [ ] Các đoạn comment không làm hỏng cú pháp Go (GoDoc syntax).
- [ ] Chạy `go build ./cmd/passiontech-beauty-contest/` thành công không có lỗi.
- [ ] Kiểm tra linter: `golangci-lint run` (nếu có cài đặt) hoặc kiểm tra cảnh báo đỏ của VS Code.
- [ ] Các comments mang tính giải thích tư duy kiến trúc (WHY and HOW), không giải thích cú pháp cơ bản (WHAT).
