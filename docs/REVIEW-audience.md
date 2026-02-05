# Báo cáo Phân tích Kỹ thuật: Module Khán giả & Hệ thống Bình chọn

Bản báo cáo này mô tả các cơ chế đảm bảo tính minh bạch và khả năng chịu tải của hệ thống bình chọn khán giả.

---

## 1. Cơ chế Chống Gian lận (Anti-Cheat Engine)

Để đảm bảo kết quả cuộc thi là công bằng nhất, chúng tôi thiết lập 3 lớp kiểm soát bình chọn:

### Lớp 1: Định danh Tài khoản (Account Limit)
Hệ thống sử dụng Redis để lưu vết bình chọn của mỗi tài khoản. Mỗi ngày, một tài khoản chỉ được phép bình chọn **01 lần** duy nhất cho mỗi thí sinh họ yêu thích.
- **Cơ chế**: Dùng khóa Cache kết hợp ngày tháng (`YYYYMMDD`) để tự động Reset giới hạn khi sang ngày mới.

### Lớp 2: Kiểm soát Địa chỉ IP (IP Rate Limiting)
Ngăn chặn việc một người dùng tạo nhiều tài khoản ảo trên cùng một thiết bị/đường truyền để bình chọn.
- **Thông minh**: Hệ thống có khả năng đọc xuyên qua các lớp Proxy/CDN để lấy IP thật của người dùng.
- **Giới hạn**: Chặn các hành vi spam từ cùng một địa chỉ IP cho cùng một thí sinh trong vòng 24 giờ.

### Lớp 3: Kiểm soát trạng thái Thí sinh
Ngăn chặn lỗi logic cho phép bình chọn cho các thí sinh chưa được duyệt hồ sơ hoặc đã bị ban khỏi cuộc thi.

---

## 2. Tối ưu Hiệu năng (High Performance)

Vì danh sách thí sinh là trang có lượt truy cập lớn nhất, chúng tôi áp dụng chiến lược Caching mạnh mẽ:
- **Cache Danh sách**: Toàn bộ danh sách hàng trăm thí sinh được lưu vào RAM (Cache TTL: 5 phút). Khi khán giả vào xem, hệ thống không cần truy vấn Database, giúp phản hồi gần như tức thì.
- **Cache Chi tiết**: Thông tin chi tiết của từng thí sinh cũng được Caching để giảm tải tối đa cho máy chủ dữ liệu.

---

## 3. Quy trình Bình chọn chuẩn (Standard Workflow)

1. **Khán giả truy cập**: Hệ thống lấy IP & User-Agent để định danh máy.
2. **Chọn thí sinh**: Kiểm tra trạng thái thí sinh (`Status == Approved`).
3. **Kiểm tra giới hạn**:
    - Kiểm tra IP trong 24h qua.
    - Kiểm tra tài khoản trong ngày hiện tại.
4. **Ghi nhận**: 
    - Tăng số lượt vote của thí sinh.
    - Ghi lại lịch sử (Audit Log) bao gồm IP, thiết bị, thời gian để đối soát nếu cần.

---

## Kết luận
Hệ thống hiện tại đã sẵn sàng cho giai đoạn bình chọn bùng nổ của khán giả với khả năng bảo mật cao và tốc độ truy cập nhanh.
