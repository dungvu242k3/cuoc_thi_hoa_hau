package handler

import (
	"cuoc_thi_hoa_hau/internal/core/domain"
	"cuoc_thi_hoa_hau/internal/core/port"
	"encoding/json"
	"net/http"
	"unicode"
)

type AuthHandler struct {
	authSvc port.AuthService
}

func NewAuthHandler(authSvc port.AuthService) *AuthHandler {
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

func (h *AuthHandler) RegisterAudience(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple Validation
	if len(req.Email) < 3 || len(req.Password) < 8 {
		http.Error(w, "Invalid email or password length", http.StatusBadRequest)
		return
	}

	// Enforce strong password
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

	claims, token, err := h.authSvc.Register(r.Context(), req.Email, req.Password, string(domain.RoleAudience))
	if err != nil {
		// Differentiate duplicate user vs other errors if possible, or just 500/409
		// AuthService returns "username already exists"
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

	// Verify Audience Role
	if claims.Role != string(domain.RoleAudience) && claims.Role != string(domain.RoleAdmin) {
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
