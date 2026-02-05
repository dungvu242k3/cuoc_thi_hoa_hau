# 👑 Beauty Contest Management System (Hệ thống Quản lý Cuộc thi Hoa hậu)

Dự án phát triển hệ thống quản lý toàn diện cho một cuộc thi sắc đẹp, bao gồm quản lý hồ sơ thí sinh, chấm điểm trực tuyến, và tương tác khán giả.

## 🚀 Tính năng nổi bật (Features)

### 🖥️ Frontend (Web App)
Xây dựng trên nền tảng hiện đại, tối ưu trải nghiệm người dùng (UX/UI).

*   **Công nghệ**: [Vite](https://vitejs.dev/) + [React](https://react.dev/) + [TypeScript](https://www.typescriptlang.org/).
*   **Styling**: [TailwindCSS v4](https://tailwindcss.com/) - Thiết kế Responsive, hiện đại.
*   **State Management**: [Zustand](https://github.com/pmndrs/zustand) - Quản lý trạng thái nhẹ nhàng, hiệu quả.
*   **Data Fetching**: [TanStack Query](https://tanstack.com/query/latest) (React Query) - Tối ưu caching và server state.
*   **Forms**: `React Hook Form` + `Zod` Validation.
*   **Charts**: `@nivo/*` - Biểu đồ thống kê trực quan.

### ⚙️ Backend (API Service)
Hệ thống API mạnh mẽ, bảo mật và hiệu năng cao.

*   **Công nghệ**: [Go](https://go.dev/) (Golang) 1.2x.
*   **API Protocol**: **GraphQL** (sử dụng thư viện `99designs/gqlgen`) - Linh hoạt trong việc truy vấn dữ liệu.
*   **Database**: **MongoDB** (NoSQL) - Lưu trữ hồ sơ linh động.
*   **Caching/Rate Limit**: **Redis** - Chống spam và tăng tốc độ phản hồi.
*   **Authentication**: JWT (JSON Web Token) + BCrypt Hashing.
*   **Architecture**: Clean Architecture / Hexagonal (Ports & Adapters).

---

## 🛠️ Cài đặt & Chạy dự án (Installation)

### 1. Yêu cầu hệ thống (Prerequisites)
*   **Go**: >= 1.20
*   **Node.js**: >= 18 (Khuyến nghị LTS)
*   **MongoDB**: Local hoặc Atlas
*   **Redis**: Local (Port 6379)

### 2. Cấu hình Backend
1.  Vào thư mục `BE`.
2.  Tạo file `.env` từ `.env.example`:
    ```bash
    cp .env.example .env
    ```
3.  Cập nhật thông tin trong `.env`:
    *   `MONGODB_URI`: Đường dẫn kết nối MongoDB.
    *   `SEED_EXAMINERS`: Cấu hình tài khoản giám khảo mặc định (JSON).

### 3. Chạy Server
**Backend**:
```bash
cd BE
go mod tidy
go run cmd/api/main.go
# Server chạy tại: http://localhost:8080
# GraphQL Playground: http://localhost:8080/
```

**Frontend**:
```bash
cd FE
npm install
npm run dev
# Web chạy tại: http://localhost:5173
```

---

## 📂 Cấu trúc dự án (Structure)

```
cuoc_thi_hoa_hau/
├── BE/                 # Backend Source Code
│   ├── cmd/            # Entry points (api, seeder)
│   ├── config/         # Configuration logic
│   ├── internal/       # Core business logic (Clean Arch)
│   │   ├── adapter/    # Handlers, Repositories, Resolvers
│   │   ├── core/       # Domain Models, Services, Ports
│   │   └── pkg/        # Shared packages (utils, validation)
│   └── docs/           # Backend Documentation
│
├── FE/                 # Frontend Source Code
│   ├── src/
│   │   ├── components/ # Reusable UI components
│   │   ├── pages/      # Route pages
│   │   ├── hooks/      # Custom React hooks
│   │   └── stores/     # Zustand stores
│   └── package.json
│
├── .agent/             # AI Agent Configurations & Workflows
└── HUONG_DAN_CHAY.md   # Quick start guide (Legacy)
```

## 🔒 Tài khoản mặc định (Seeder)
Hệ thống hỗ trợ seed dữ liệu mẫu cho Permissions, Roles.

```bash
cd BE
go run cmd/seeder/main.go
```
