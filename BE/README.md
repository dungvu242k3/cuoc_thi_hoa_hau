# Backend Cuộc Thi Hoa Hậu (Miss Beauty Contest)

## 📌 Tổng Quan
Backend Server quản lý toàn diện hệ thống cuộc thi, viết bằng **Go (Golang)** + **GraphQL** + **MongoDB**.
Kiến trúc: **Clean Architecture** (Domain -> Port -> Service -> Adapter).

---

## 🔥 Tính Năng Đã Triển Khai

| Module | Chức Năng | Trạng Thái | Mô Tả |
| :--- | :--- | :--- | :--- |
| **1. Authentication** | Đăng ký, Đăng nhập (JWT) | ✅ Hoàn tất | Bảo mật JWT, phân quyền Role (CANDIDATE, ADMIN...) |
| **2. Contestant** | Quản lý Profile thí sinh | ✅ Hoàn tất | Đầy đủ vòng đời (Draft -> Pending -> Approved/Rejected). Bảo vệ PII 3 lớp. |
| **3. Schedule** | Lịch trình & Sự kiện | ✅ Hoàn tất | Lịch thi đấu, tập luyện. Phân trang. |
| **4. Feedback** | Gửi khiếu nại/góp ý | ✅ Hoàn tất | Thí sinh gửi yêu cầu, Admin xử lý. |
| **5. Score** | Xem điểm thi | ✅ Hoàn tất | Xem điểm chi tiết từng phần thi. Bảo mật chỉ xem điểm chính chủ. |

---

## 🛠️ Chi Tiết Chức Năng & Hàm Quan Trọng (Internals)

### 1. Module Hồ Sơ (Contestant)
*   `CreateProfile`: Tạo mới hồ sơ (Draft). Kiểm tra trùng lặp UserID/CCCD. Valid tuổi > 18, chiều cao > 1m60.
*   `UpdateProfile`: Cập nhật thông tin. Chỉ cho phép khi trạng thái là DRAFT hoặc REJECTED. Tự động log audit.
*   `SubmitProfile`: Chốt hồ sơ, chuyển sang PENDING. Sau bước này thí sinh không thể sửa đổi.
*   `ToPublicView`: **(Bảo mật)** Hàm domain chủ động xóa các trường nhạy cảm (SĐT, Email, Địa chỉ) trước khi trả về API Public.

### 2. Module Phản Hồi (Feedback)
*   `SendFeedback`: Gửi feedback loại `PROPOSAL`, `COMPLAINT`, hoặc `REQUEST`. Tự động gán trạng thái `PENDING`.

### 3. Module Điểm Số (Score)
*   `GetMyScore`: Lấy điểm của thí sinh đang đăng nhập.
    *   **Bảo mật**: Xác thực UserID từ Token JWT -> Tra cứu ra CandidateID -> Lấy điểm. Ngăn chặn tuyệt đối việc xem trộm điểm người khác.

---

## 📡 API Documentation (GraphQL)

### A. Nhóm Thí Sinh (Yêu cầu đăng nhập)

#### 1. Hồ sơ cá nhân
```graphql
query {
  myProfile {
    id
    status
    personalInfo { fullName phone email address }
    physicalInfo { height weight measurements }
    portfolio { introduction galleryUrls }
  }
}
```

#### 2. Thao tác hồ sơ
```graphql
# Tạo mới
mutation {
  createContestantProfile(input: {
    fullName: "Nguyen Van A",
    identifyCard: "0123456789",
    # ...
  }) { id }
}

# Cập nhật
mutation {
  updateContestantProfile(input: { height: 172.5 }) { id }
}

# Nộp hồ sơ (Chốt)
mutation { submitProfile }

# Xóa hồ sơ (Rút lui)
mutation { deleteProfile }
```

#### 3. Gửi Phản Hồi / Khiếu Nại
```graphql
mutation {
  sendFeedback(input: {
    title: "Sai sót điểm số",
    content: "Tôi thấy điểm phần thi áo tắm chưa chính xác...",
    type: COMPLAINT 
  })
}
```

#### 4. Xem danh sách Phản hồi của tôi
```graphql
query {
  myFeedbacks(limit: 10, offset: 0) {
    items { title status createdAt }
    total
  }
}
```

#### 5. Xem Điểm Số (Mới ⭐️)
```graphql
query {
  myScore {
    totalScore
    details {
      key    # Tên tiêu chí (VD: "Hình thể")
      value  # Điểm số (VD: 9.5)
    }
  }
}
```

---

### B. Nhóm Cộng Đồng & Khán Giả (Community) 

#### 1. Danh sách thí sinh & Chi tiết
*Dữ liệu trả về đã được ẩn thông tin nhạy cảm.*
```graphql
query {
  publicContestants(limit: 10, page: 1) {
    id
    personalInfo { fullName }
    portfolio { avatarUrl }
  }
}

```

#### 2. Lịch trình cuộc thi
```graphql
query {
  publicSchedules(limit: 5) {
    items { title startTime location type }
  }
}
```
