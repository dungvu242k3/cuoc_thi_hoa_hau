/*
Package ecode (Error Codes) quản lý toàn bộ các mã lỗi Business (nghiệp vụ) của toàn dự án.
Tác dụng:
- Ngăn chặn lỗi "Magic String", nơi mỗi người code tự gõ message lỗi theo ý mình.
- Tầng Service chỉ cần: return ecode.ErrProfileExists, Middleware hoặc Handler sẽ tự biết map sang HTTP 400.
*/
package ecode

import "errors"

// Auth errors
var (
	ErrUnauthorized    = errors.New("unauthorized")
	ErrInvalidCredent  = errors.New("invalid credentials")
	ErrUserExists      = errors.New("username already exists")
	ErrForbidden       = errors.New("bạn không có quyền thực hiện hành động này")
	ErrTooManyAttempts = errors.New("too many attempts, please try again later")
)

// Contestant errors
var (
	ErrProfileExists   = errors.New("bạn đã có hồ sơ, không thể tạo thêm")
	ErrProfileNotFound = errors.New("hồ sơ không tồn tại")
	ErrDuplicateIDCard = errors.New("số CCCD này đã được sử dụng")
	ErrProfileLocked   = errors.New("hồ sơ đã xử lý, không được phép chỉnh sửa")
	ErrAvatarRequired  = errors.New("vui lòng cập nhật ảnh đại diện")
	ErrNotApproved     = errors.New("chỉ có thể bình chọn cho thí sinh đã qua vòng duyệt hồ sơ")
	ErrContestantGone  = errors.New("thí sinh không tồn tại")
)

// Vote errors
var (
	ErrAlreadyVoted   = errors.New("bạn đã bình chọn cho thí sinh này rồi")
	ErrIPLimitReached = errors.New("địa chỉ IP này đã bình chọn cho thí sinh này rồi")
	ErrVoteDailyLimit = errors.New("bạn đã bình chọn cho thí sinh này trong hôm nay rồi")
	ErrSystemBusy     = errors.New("hệ thống đang bận, vui lòng thử lại sau")
)

// Score errors
var (
	ErrMissingContestant = errors.New("thiếu ID thí sinh")
	ErrNoScoreCriteria   = errors.New("vui lòng nhập điểm cho ít nhất một tiêu chí")
)
