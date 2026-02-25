/*
Package validation bao gồm các tiện ích làm sạch (sanitize) dữ liệu đầu vào.
Tác dụng của pagination.go:
- Ngăn chặn Frontend truyền `limit = 999999` làm sập Database (DoS).
- Ép cứng mức tối đa (MaxPageSize) và mặc định (DefaultPageSize) cho mọi API phân trang.
*/
package validation

const (
	MaxPageSize     = 100
	DefaultPageSize = 20
)

// NormalizePagination chuẩn hóa các tham số phân trang, fallback về default nếu truyền sai hoặc truyền nil.
func NormalizePagination(page *int, limit *int) (int, int) {
	p := 1
	l := DefaultPageSize

	if page != nil && *page > 0 {
		p = *page
	}

	if limit != nil && *limit > 0 {
		l = *limit
		if l > MaxPageSize {
			l = MaxPageSize
		}
	}

	return p, l
}
