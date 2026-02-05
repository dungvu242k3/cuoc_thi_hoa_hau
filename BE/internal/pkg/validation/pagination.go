package validation

const (
	MaxPageSize     = 100
	DefaultPageSize = 20
)

// NormalizePagination ensures safe pagination parameters
func NormalizePagination(page *int, limit *int) (int, int) {
	p := 1
	l := DefaultPageSize

	if page != nil && *page > 0 {
		p = *page
	}

	if limit != nil && *limit > 0 {
		l = *limit
		// Enforce max limit to prevent DoS
		if l > MaxPageSize {
			l = MaxPageSize
		}
	}

	return p, l
}
