# Báo cáo Phân tích Kỹ thuật: Module Ban Giám Khảo (Examiner Scoring)

Bản báo cáo này giải thích các cơ chế hoạt động, tính bảo mật và tối ưu hóa của hệ thống chấm điểm dành cho Ban giám khảo.

---

## 1. Cơ chế Chấm điểm & Tối ưu hóa

### Tự động tính tổng (Auto-Calculation)
Hệ thống được thiết kế để Ban giám khảo chỉ cần nhập điểm cho từng tiêu chí (như Gương mặt, Hình thể, Kỹ năng ứng xử). Hệ thống Backend sẽ tự động tính tổng điểm (`TotalScore`) một cách chính xác trước khi lưu vào cơ sở dữ liệu.

### Tối ưu hóa gán biến (Mappers)
Chúng tôi áp dụng mô hình chuyên nghiệp để tách biệt dữ liệu Giao diện (GraphQL) và dữ liệu Nghiệp vụ (Domain).
- **Lợi ích**: Giúp mã nguồn cực kỳ dễ đọc. Mọi thay đổi về cấu trúc dữ liệu trong tương lai sẽ chỉ cần chỉnh sửa tại một nơi duy nhất (`mapper.go`), không làm ảnh hưởng đến logic xử lý chính.

---

## 2. Các lớp Bảo mật (Security Layers)

Hệ thống chấm điểm là phần nhạy cảm nhất, do đó chúng tôi áp dụng 4 lớp bảo vệ:

### Lớp 1: Xác thực & Phân quyền (RBAC)
Chỉ những người dùng có quyền `score:write` (được cấp cho tài khoản Ban giám khảo) mới có thể gọi hàm chấm điểm. Hệ thống tự động chặn mọi nỗ lực truy cập trái phép ngay từ tầng Gateway.

### Lớp 2: Kiểm soát dữ liệu (Validation)
- **Ràng buộc giá trị**: Điểm số được kiểm soát chặt chẽ trong khung từ 0 đến 10. Mọi giá trị nằm ngoài khoảng này sẽ bị từ chối ngay lập tức.
- **Tính toàn vẹn**: Hệ thống kiểm tra sự tồn tại của Thí sinh trước khi nhận điểm, đảm bảo không có điểm "ma" được ghi vào hệ thống.

### Lớp 3: Chống tấn công nội dung (XSS)
Tất cả các lời nhận xét (Comment) của giám khảo đều được đi qua bộ lọc **HTML Escaping**. Điều này ngăn chặn việc giám khảo (hoặc người giả mạo) chèn các đoạn mã script độc hại vào hệ thống để tấn công người xem hoặc quản trị viên.

### Lớp 4: Cơ chế Audit Log (Truy vết)
Đây là tính năng quan trọng nhất cho tính minh bạch:
- **Ghi lại IP & Thiết bị**: Mỗi lần chấm điểm, hệ thống sẽ lưu lại địa chỉ IP mạng và thông tin trình duyệt của giám khảo.
- **Lợi ích**: Nếu có bất kỳ nghi ngờ nào về việc gian lận hoặc tài khoản bị chiếm đoạt, Ban tổ chức có thể dựa vào thông tin này để đối soát và xác minh.

---

## 3. Quy trình hoạt động (Workflow)

1. **Giám khảo đăng nhập**: Hệ thống cấp Token định danh.
2. **Nhập điểm**: Giám khảo chọn thí sinh và nhập điểm + nhận xét.
3. **Backend xử lý**:
    - Kiểm tra quyền hạn.
    - Validate điểm số (0-10).
    - Lọc XSS cho nhận xét.
    - **Ghi nhận IP/Thiết bị**.
    - Tính tổng điểm.
4. **Lưu trữ**: Dữ liệu được lưu an toàn xuống MongoDB.

---

## Kết luận
Hệ thống chấm điểm hiện tại đã đạt tiêu chuẩn về **Tính minh bạch** và **Tính bảo mật cao**, sẵn sàng phục vụ cho các vòng thi quan trọng của cuộc thi.
