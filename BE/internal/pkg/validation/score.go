package validation

import (
	"errors"
	"fmt"
)

// ValidateScore checks if the criteria scores are within valid range (0-10)
func ValidateScore(criteriaScores map[string]float64) error {
	if len(criteriaScores) == 0 {
		return errors.New("vui lòng nhập điểm cho ít nhất một tiêu chí")
	}

	for key, score := range criteriaScores {
		if score < 0 || score > 10 {
			return fmt.Errorf("điểm cho tiêu chí '%s' phải nằm trong khoảng từ 0 đến 10", key)
		}
	}

	return nil
}
