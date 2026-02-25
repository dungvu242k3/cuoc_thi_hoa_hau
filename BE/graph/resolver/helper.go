package resolver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"
	"unicode"
)

// validateAuthInput validates email and password for auth operations
func validateAuthInput(email, password string) error {
	if len(email) < 3 {
		return errors.New("email must be at least 3 characters")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least one uppercase, lowercase, number, and special character")
	}

	return nil
}

// checkRateLimit implements per-action rate limiting using cache
func (r *mutationResolver) checkRateLimit(ctx context.Context, action, key string) error {
	if r.CacheSvc == nil {
		return nil
	}

	limitKey := fmt.Sprintf("rate_limit:%s:%s", action, key)
	count, err := r.CacheSvc.Incr(ctx, limitKey)
	if err != nil {
		slog.Error("Rate limit counter failed", "error", err, "key", limitKey)
		return nil
	}

	if count == 1 {
		r.CacheSvc.Set(ctx, limitKey, 1, 1*time.Minute)
	}

	if count > 5 {
		return errors.New("too many attempts, please try again later")
	}
	return nil
}

// checkVoteRateLimit implements per-user per-contestant daily rate limiting
func (r *mutationResolver) checkVoteRateLimit(ctx context.Context, userID, contestantID string) error {
	date := time.Now().Format("20060102")
	key := fmt.Sprintf("vote_limit:%s:%s:%s", date, userID, contestantID)

	count, err := r.CacheSvc.Incr(ctx, key)
	if err != nil {
		log.Printf("[SECURITY] Rate limit cache error: %v", err)
		return errors.New("hệ thống đang bận, vui lòng thử lại sau")
	}

	if count == 1 {
		_ = r.CacheSvc.Set(ctx, key, "1", 24*time.Hour)
	}

	if count > 1 {
		return errors.New("bạn đã bình chọn cho thí sinh này trong hôm nay rồi")
	}

	return nil
}

func safelyToFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case json.Number:
		return val.Float64()
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case float32:
		return float64(val), nil
	default:
		return 0, fmt.Errorf("unknown type")
	}
}
