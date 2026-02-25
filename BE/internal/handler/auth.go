/*
Package handler đóng vai trò là Controller trong mô hình MVC.
Tác dụng:
- Là tầng đầu tiên ngoài cùng tiếp xúc với Client (Frontend/Mobile).
- Nhận Request (JSON), kiểm tra lỗi cơ bản (Validation), sau đó gọi Tầng Service xử lý nghiệp vụ.
- Cuối cùng, gói kết quả từ Service thành JSON trả về cho Client.
- KHÔNG BAO GIỜ viết logic truy vấn DB hay logic kinh doanh phức tạp ở đây.
*/
package handler

import (
	"cuoc_thi_hoa_hau/internal/model"
	"cuoc_thi_hoa_hau/internal/types"
	"encoding/json"
	"net/http"
	"unicode"
)

type AuthHandler struct {
	authSvc types.AuthService
}

func NewAuthHandler(authSvc types.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
	Role   string `json:"role"`
}

// RegisterAudience xử lý luồng đăng ký tài khoản cho Khán giả.
// Nó parse chuỗi JSON từ body HTTP request thành dạng Struct của Go.
func (h *AuthHandler) RegisterAudience(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Email) < 3 || len(req.Password) < 8 {
		http.Error(w, "Invalid email or password length", http.StatusBadRequest)
		return
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)
	for _, char := range req.Password {
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
		http.Error(w, "Password must contain uppercase, lowercase, number, and special char", http.StatusBadRequest)
		return
	}

	// Gọi xuống tầng Service (Tầng nghiệp vụ) để thực thi logic tạo user và sinh token.
	// Bất cứ nghiệp vụ lõi nào cũng nằm ở Service, Handler chỉ việc đứng đợi kết quả.
	claims, token, err := h.authSvc.Register(r.Context(), req.Email, req.Password, string(model.RoleAudience))
	if err != nil {
		if err.Error() == "username already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := AuthResponse{
		Token:  token,
		UserID: claims.UserID,
		Role:   claims.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) LoginAudience(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, token, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if claims.Role != string(model.RoleAudience) && claims.Role != string(model.RoleAdmin) {
		http.Error(w, "Unauthorized access for this portal", http.StatusForbidden)
		return
	}

	resp := AuthResponse{
		Token:  token,
		UserID: claims.UserID,
		Role:   claims.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
