/*
Package validation chứa các hàm kiểm tra tính toàn vẹn của dữ liệu đầu vào.
Tác dụng của contestant.go:
- Lọc đầu vào cực kì nghiêm ngặt trước khi dội xuống Service và DB.
- Chặn các nguy cơ XSS (như HTML tags: <script>), kiểm tra độ dài, regex email, sđt hợp lệ.
*/
package validation

import (
	gqlmodel "cuoc_thi_hoa_hau/graph/model"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ... (Rest of the file remains exactly the same, I will use write_to_file with full contents)
// ValidateCreateContestantInput là lưới lọc đầu tiên khi Frontend gửi Mutation tạo hồ sơ.
// Gọi các hàm kiểm tra con (ValidateName, ValidateEmail...) cho từng field.
func ValidateCreateContestantInput(input gqlmodel.CreateContestantInput) error {
	if err := ValidateName(input.FullName); err != nil {
		return fmt.Errorf("fullName: %w", err)
	}

	if err := ValidateEmail(input.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}

	if err := ValidatePhone(input.Phone); err != nil {
		return fmt.Errorf("phone: %w", err)
	}

	dob, err := time.Parse("2006-01-02", input.Dob)
	if err != nil {
		return errors.New("dob: invalid date format, expected YYYY-MM-DD")
	}
	if err := ValidateDOB(dob); err != nil {
		return fmt.Errorf("dob: %w", err)
	}

	if input.Height <= 0 || input.Height > 300 {
		return errors.New("height: must be between 1-300 cm")
	}
	if input.Weight <= 0 || input.Weight > 500 {
		return errors.New("weight: must be between 1-500 kg")
	}

	if input.AvatarURL != nil {
		if err := ValidateURL(*input.AvatarURL); err != nil {
			return fmt.Errorf("avatarURL: %w", err)
		}
	}

	if input.GalleryUrls != nil {
		for i, galleryURL := range input.GalleryUrls {
			if err := ValidateURL(galleryURL); err != nil {
				return fmt.Errorf("galleryUrls[%d]: %w", i, err)
			}
		}
	}

	if input.SocialLinks != nil {
		for platform, link := range input.SocialLinks {
			if err := ValidateURL(link); err != nil {
				return fmt.Errorf("socialLinks[%v]: %w", platform, err)
			}
		}
	}

	if input.Introduction != nil {
		if err := ValidateText(*input.Introduction, 5000); err != nil {
			return fmt.Errorf("introduction: %w", err)
		}
	}

	if err := ValidateText(input.Address, 500); err != nil {
		return fmt.Errorf("address: %w", err)
	}

	if err := ValidateText(input.Job, 200); err != nil {
		return fmt.Errorf("job: %w", err)
	}

	if err := ValidateText(input.Nationality, 100); err != nil {
		return fmt.Errorf("nationality: %w", err)
	}

	return nil
}

// ValidateUpdateContestantInput cho phép truyền các field nil (vì graphql cho phép người dùng muốn sửa field nào thì gửi field đó).
func ValidateUpdateContestantInput(input gqlmodel.UpdateContestantInput) error {
	if input.FullName != nil {
		if err := ValidateName(*input.FullName); err != nil {
			return fmt.Errorf("fullName: %w", err)
		}
	}

	if input.Email != nil {
		if err := ValidateEmail(*input.Email); err != nil {
			return fmt.Errorf("email: %w", err)
		}
	}

	if input.Phone != nil {
		if err := ValidatePhone(*input.Phone); err != nil {
			return fmt.Errorf("phone: %w", err)
		}
	}

	if input.Dob != nil {
		dob, err := time.Parse("2006-01-02", *input.Dob)
		if err != nil {
			return errors.New("dob: invalid date format, expected YYYY-MM-DD")
		}
		if err := ValidateDOB(dob); err != nil {
			return fmt.Errorf("dob: %w", err)
		}
	}

	if input.Height != nil {
		if *input.Height <= 0 || *input.Height > 300 {
			return errors.New("height: must be between 1-300 cm")
		}
	}

	if input.Weight != nil {
		if *input.Weight <= 0 || *input.Weight > 500 {
			return errors.New("weight: must be between 1-500 kg")
		}
	}

	if input.AvatarURL != nil {
		if err := ValidateURL(*input.AvatarURL); err != nil {
			return fmt.Errorf("avatarURL: %w", err)
		}
	}

	if input.GalleryUrls != nil {
		for i, galleryURL := range input.GalleryUrls {
			if err := ValidateURL(galleryURL); err != nil {
				return fmt.Errorf("galleryUrls[%d]: %w", i, err)
			}
		}
	}

	if input.SocialLinks != nil {
		for platform, link := range input.SocialLinks {
			if err := ValidateURL(link); err != nil {
				return fmt.Errorf("socialLinks[%v]: %w", platform, err)
			}
		}
	}

	if input.Introduction != nil {
		if err := ValidateText(*input.Introduction, 5000); err != nil {
			return fmt.Errorf("introduction: %w", err)
		}
	}

	if input.Address != nil {
		if err := ValidateText(*input.Address, 500); err != nil {
			return fmt.Errorf("address: %w", err)
		}
	}

	if input.Job != nil {
		if err := ValidateText(*input.Job, 200); err != nil {
			return fmt.Errorf("job: %w", err)
		}
	}

	if input.Nationality != nil {
		if err := ValidateText(*input.Nationality, 100); err != nil {
			return fmt.Errorf("nationality: %w", err)
		}
	}

	return nil
}

func ValidateName(name string) error {
	if len(name) < 2 || len(name) > 100 {
		return errors.New("must be between 2-100 characters")
	}

	if ContainsHTMLTags(name) {
		return errors.New("cannot contain HTML tags")
	}

	return nil
}

func ValidateEmail(email string) error {
	if len(email) > 254 {
		return errors.New("email too long")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePhone(phone string) error {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(phone, "-", ""), " ", "")
	cleaned = strings.ReplaceAll(cleaned, "(", "")
	cleaned = strings.ReplaceAll(cleaned, ")", "")

	if len(cleaned) < 10 || len(cleaned) > 15 {
		return errors.New("must be between 10-15 digits")
	}

	phoneRegex := regexp.MustCompile(`^[+]?[0-9]{10,15}$`)
	if !phoneRegex.MatchString(cleaned) {
		return errors.New("invalid phone format")
	}
	return nil
}

func ValidateDOB(dob time.Time) error {
	now := time.Now()
	if dob.After(now) {
		return errors.New("cannot be in the future")
	}

	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}

	if age < 18 {
		return errors.New("contestant must be at least 18 years old")
	}
	if age > 100 {
		return errors.New("invalid age")
	}

	return nil
}

func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return nil
	}

	if len(urlStr) > 2048 {
		return errors.New("URL too long")
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return errors.New("malformed URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("must use http or https protocol")
	}

	return nil
}

func ValidateText(text string, maxLength int) error {
	if utf8.RuneCountInString(text) > maxLength {
		return fmt.Errorf("exceeds maximum length of %d characters", maxLength)
	}

	if ContainsHTMLTags(text) {
		return errors.New("cannot contain HTML/script tags")
	}

	return nil
}

func ContainsHTMLTags(s string) bool {
	if s == "" {
		return false
	}

	lower := strings.ToLower(s)

	tagRegex := regexp.MustCompile(`<[a-zA-Z/!].*?>`)
	if tagRegex.MatchString(s) {
		return true
	}

	dangerousPatterns := []string{
		"javascript:",
		"vbscript:",
		"data:text/html",
		"onerror=",
		"onclick=",
		"onload=",
		"onmouseover=",
		"onfocus=",
		"onabort=",
		"onkeydown=",
		"onkeyup=",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}
