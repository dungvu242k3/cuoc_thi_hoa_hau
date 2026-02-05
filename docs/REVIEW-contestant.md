# Báo cáo Rà soát & Phân tích Kỹ thuật: Module Thí sinh (Contestant)

Bản báo cáo này cung cấp cái nhìn tổng quan về cách module Thí sinh được xây dựng, tối ưu hóa và bảo vệ trong hệ thống.

---

## 1. Kiến trúc Hàm & Biến (Optimization)

### Cơ chế Mapper (Khử mã lặp)
Chúng tôi sử dụng mô hình **DTO -> Domain Mapping**. Thay vì để logic chuyển đổi dữ liệu nằm rải rác trong các hàm điều hướng (Resolvers), toàn bộ logic này được tập trung tại `resolver/mapper.go`.

- **Lợi ích**:
    - **Sạch sẽ**: Hàm Resolver giờ đây chỉ còn khoảng 10 dòng thay vì 50-60 dòng như trước.
    - **An toàn**: Sử dụng các hàm helper như `getString` để xử lý con trỏ (nil pointer), tránh tình trạng ứng dụng bị "crash" khi dữ liệu đầu vào thiếu.
    - **Tối ưu biến**: Giảm thiểu việc khởi tạo các biến tạm thời không cần thiết, giúp bộ nhớ (RAM) được sử dụng hiệu quả hơn.

---

## 2. Cơ chế Bảo mật (Security)

Hệ thống được thiết kế với "Bảo mật đa tầng" (Defense in Depth):

### A. Chống tấn công XSS (Cross-Site Scripting)
Dữ liệu từ người dùng (như phần giới thiệu) luôn tiềm ẩn mã độc. Chúng tôi sử dụng bộ lọc kết hợp:
- **Regex Detection**: Tự động phát hiện và chặn mọi thẻ HTML dạng `<tag>`.
- **Event Handler Blacklist**: Chặn các thuộc tính nguy hiểm như `onclick`, `onerror`, `javascript:`, giúp ngăn chặn việc thực thi mã JavaScript lạ trong trình duyệt của người xem.

### B. Kiểm soát Bình chọn (Rate Limiting)
Để đảm bảo tính công minh, chức năng bình chọn được trang bị:
- **Atomic Operation**: Sử dụng lệnh `Incr` (Atomic) để đảm bảo dù có hàng ngàn người nhấn cùng lúc, số liệu vẫn chính xác 100%.
- **Fail-Closed Policy**: Nếu hệ thống kiểm soát (Redis/Cache) gặp sự cố, hệ thống sẽ **từ chối** bình chọn thay vì để người dùng bình chọn tự do. Đây là tiêu chuẩn cho các cuộc thi chuyên nghiệp.

### C. Phân quyền & Bảo mật thông tin (RBAC & Masking)
Hệ thống tự động điều chỉnh dữ liệu trả về tùy theo vai trò:
- **Công chúng**: Chỉ thấy ảnh, tên, SBD. Các thông tin nhạy cảm (Số CCCD, SĐT, Email) sẽ bị che (VD: `098******`).
- **Thí sinh (Owner)**: Thấy đầy đủ thông tin của chính mình.
- **Ban tổ chức (Admin)**: Thấy toàn bộ hồ sơ để phục vụ việc xét duyệt.

---

## 3. Quy trình Chức năng (Workflow)

1. **Đăng ký**: Validate dữ liệu -> Map sang Domain -> Lưu DB -> Tạo SBD tự động.
2. **Xét duyệt**: Ban tổ chức kiểm tra hồ sơ. Khi được duyệt, thí sinh mới xuất hiện trên trang chủ.
3. **Bình chọn**: Kiểm tra IP -> Kiểm tra giới hạn 24h -> Ghi nhận lịch sử (Audit Log) -> Tăng điểm số.

---

## Kết luận
Module Thí sinh không chỉ tập trung vào tính năng mà còn đặt sự **An toàn dữ liệu** và **Công bằng hệ thống** lên hàng đầu. Các cải tiến về Mapper giúp mã nguồn đạt tiêu chuẩn **Clean Code**, sẵn sàng cho việc mở rộng quy mô lớn.
