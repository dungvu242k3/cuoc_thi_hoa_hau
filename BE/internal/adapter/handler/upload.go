package handler

import (
	"crypto/rand"
	"cuoc_thi_hoa_hau/internal/core/port"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
)

// MAX_UPLOAD_SIZE is 10MB
const MAX_UPLOAD_SIZE = 10 << 20

type FileHandler struct {
	Storage port.FileStorage
}

func NewFileHandler(storage port.FileStorage) *FileHandler {
	return &FileHandler{
		Storage: storage,
	}
}

// UploadFile handles file uploads authentication and validation, then delegates storage
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce Max Upload Size
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "File too large (Max 10MB)", http.StatusBadRequest)
		return
	}

	// 2. Retrieve the file
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3. MIME Sniffing
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(buff)
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !validTypes[contentType] {
		http.Error(w, "Invalid file type. Only JPEG, PNG, and WebP are allowed.", http.StatusBadRequest)
		return
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// 4. Create Safe Filename
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		http.Error(w, "Invalid file extension", http.StatusBadRequest)
		return
	}

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	filename := hex.EncodeToString(randomBytes) + ext

	// 5. Use Storage Interface to Save
	// The handler doesn't know (or care) if it's Local disk or S3
	url, err := h.Storage.Upload(r.Context(), file, filename)
	if err != nil {
		http.Error(w, "Storage Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Return the URL
	response := map[string]string{
		"url": url,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
