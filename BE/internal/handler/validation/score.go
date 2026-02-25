/*
Package validation kiểm tra tính hợp lệ của Phiếu Điểm.
Tác dụng của score.go:
- Chặn đứng các trường hợp cố tình đục khoét truyền điểm âm (-5.0) hoặc điểm siêu thực (99.9) vào DB.
*/
package validation

import (
	"errors"
	"fmt"
)

// ValidateScore rà soát từng tiêu chí trong bảng điểm, đảm bảo mọi điểm số đều nằm trong thang [0.0 - 10.0].
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
