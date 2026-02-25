# Kế hoạch Chú thích TOÀN BỘ Code (Annotate Full BE Codebase)

## 1. Overview (Tổng quan)
- **Mục tiêu**: Bơm chú thích (comments) giải thích kiến trúc, logic luồng, chức năng hàm vào **TOÀN BỘ** các file mã nguồn `.go` bên trong dự án Backend (`c:\Users\dungv\cuoc_thi_hoa_hau\BE`). 
- **Phạm vi Dứt khoát**: Không bỏ sót một file nào. Đi từ thư mục `cmd/`, `internal/` cho tới tận `graph/`. Giúp code như một cuốn sách giáo khoa định hướng kiến trúc ứng dụng chuẩn.

## 2. Project Type
- **BACKEND** (Golang, MongoDB, GraphQL, REST, Redis)

## 3. Success Criteria (Tiêu chí thành công)
- Bất kì file `.go` nào mở lên cũng sẽ thấy GoDoc comment trên định nghĩa Package (`package xxx`), định nghĩa Interface/Struct, và định nghĩa của các Functional logic.
- Giải thích **WHY** (tại sao phải viết hàm này ở tầng này) & **HOW** (hàm này đang làm gì).
- Code phải biên dịch thành công 100% sau khi sửa.

## 4. Tech Stack
- Go 1.24, gqlgen (GraphQL), go-chi (REST/Router), MongoDB driver, JWT, Redis.

## 5. Danh Sách 100% Các File Cần Xử Lý

Kế hoạch này sẽ duyệt toàn bộ thư mục một cách hệ thống (Systematic Sweep):

*Hệ thống thư mục lõi HTTP Server & Routing*
- [ ] `internal/server/http.go`, `bcrypt.go`, `jwt.go`
- [ ] `internal/routers/router.go`
- [ ] `internal/handler/auth.go`, `upload.go`
- [ ] `internal/handler/validation/contestant.go`, `pagination.go`, `score.go`
- [ ] `internal/config/config.go`
- [ ] `internal/ecode/ecode.go`

*Hệ thống Model & DTOs*
- [ ] `internal/model/contestant.go`, `user.go`, `score.go`, `feedback.go`, `schedule.go`, `constants.go`
- [ ] `internal/types/contestant.go`, `user.go`, `auth.go`, `score.go`, `schedule.go`, `feedback.go`, `cache.go`, `security.go`, `upload.go`

*Hệ thống Service (Business Logic Core)*
- [ ] `internal/service/auth.go`, `contestant.go`, `score.go`, `schedule.go`, `feedback.go`

*Hệ thống DAO (Database Storage)*
- [ ] `internal/dao/contestant.go`, `user.go`, `score.go`, `schedule.go`, `feedback.go`, `storage.go`
- [ ] `internal/database/connection.go`
- [ ] `internal/cache/redis_adapter.go`

*Hệ thống GraphQL*
- [ ] `graph/resolver/schema.resolvers.go`, `auth.resolvers.go`, `contestant.resolvers.go`, `scoring.resolvers.go`, `schedule.resolvers.go`, `feedback.resolvers.go`
- [ ] `graph/resolver/resolver.go`, `mapper.go`
- [ ] `graph/middleware/auth.go`, `client_info.go`, `ratelimit.go`

## 6. Execution Strategy
Toàn bộ dự án có khoảng ~40 files. Sẽ được chạy theo từng cụm Layer từ dưới lên trên:
- **Phase A**: Annotate Tầng Hạ Tầng (DB, Cache, Config, Server)
- **Phase B**: Annotate Tầng Giao Tiếp Hợp Đồng (Types, Model, Mapper, Validation)
- **Phase C**: Annotate Tầng Lưu Trữ (DAO Repo)
- **Phase D**: Annotate Tầng Nghiệp Vụ Chuyên Kĩ Thuật (Service)
- **Phase E**: Annotate Tầng Mạng Tiếp Xúc Client (Handler, Routers, Middleware, Resolvers)

## 7. Phase X: Verification
- [ ] Lệnh `go build ./cmd/passiontech-beauty-contest/` chạy pass không bị dội bom cú pháp.
- [ ] Mở ngẫu nhiên bất kỳ file nào cũng phải thấy Annotation ở dạng chuẩn.
