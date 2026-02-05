package validation

import (
	"cuoc_thi_hoa_hau/internal/adapter/graph/model"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ValidateCreateContestantInput validates all contestant input fields
func ValidateCreateContestantInput(input model.CreateContestantInput) error {
	// Validate Full Name
	if err := ValidateName(input.FullName); err != nil {
		return fmt.Errorf("fullName: %w", err)
	}

	// Validate Email
	if err := ValidateEmail(input.Email); err != nil {
		return fmt.Errorf("email: %w", err)
	}

	// Validate Phone
	if err := ValidatePhone(input.Phone); err != nil {
		return fmt.Errorf("phone: %w", err)
	}

	// Validate DOB
	dob, err := time.Parse("2006-01-02", input.Dob)
	if err != nil {
		return errors.New("dob: invalid date format, expected YYYY-MM-DD")
	}
	if err := ValidateDOB(dob); err != nil {
		return fmt.Errorf("dob: %w", err)
	}

	// Validate Physical Measurements
	if input.Height <= 0 || input.Height > 300 {
		return errors.New("height: must be between 1-300 cm")
	}
	if input.Weight <= 0 || input.Weight > 500 {
		return errors.New("weight: must be between 1-500 kg")
	}

	// Validate URLs
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

	// Validate Text Fields
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

// ValidateUpdateContestantInput validates update input fields
func ValidateUpdateContestantInput(input model.UpdateContestantInput) error {
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

// ValidateName ensures name contains no scripts or excessive length
func ValidateName(name string) error {
	if len(name) < 2 || len(name) > 100 {
		return errors.New("must be between 2-100 characters")
	}

	if ContainsHTMLTags(name) {
		return errors.New("cannot contain HTML tags")
	}

	return nil
}

// ValidateEmail validates email format
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

// ValidatePhone validates phone number
func ValidatePhone(phone string) error {
	// Remove common separators
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

// ValidateDOB ensures date of birth is in the past and reasonable
func ValidateDOB(dob time.Time) error {
	now := time.Now()
	if dob.After(now) {
		return errors.New("cannot be in the future")
	}

	// Age must be between 18-100 years
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

// ValidateURL ensures URL is safe (http/https only, no javascript:)
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

	// Only allow http/https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("must use http or https protocol")
	}

	return nil
}

// ValidateText checks for XSS and length limits
func ValidateText(text string, maxLength int) error {
	if utf8.RuneCountInString(text) > maxLength {
		return fmt.Errorf("exceeds maximum length of %d characters", maxLength)
	}

	if ContainsHTMLTags(text) {
		return errors.New("cannot contain HTML/script tags")
	}

	return nil
}

// ContainsHTMLTags detects common HTML/script tags and structural patterns
func ContainsHTMLTags(s string) bool {
	if s == "" {
		return false
	}

	lower := strings.ToLower(s)

	// 1. Regex check for Any HTML Tag: <tag ...> or </tag>
	// This covers most standard and non-standard tags
	tagRegex := regexp.MustCompile(`<[a-zA-Z/!].*?>`)
	if tagRegex.MatchString(s) {
		return true
	}

	// 2. Blacklist check for event handlers and dangerous protocols
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
