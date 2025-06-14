package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ufo/internal/config"
)

// UploadHandler handles file upload operations
type UploadHandler struct {
	config   *config.Config
	template *template.Template
}

// FileInfo represents uploaded file information
type FileInfo struct {
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	UploadTime time.Time `json:"upload_time"`
}

// UploadResponse represents the JSON response for uploads
type UploadResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	File    *FileInfo `json:"file,omitempty"`
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(cfg *config.Config) *UploadHandler {
	// Parse HTML template
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	return &UploadHandler{
		config:   cfg,
		template: tmpl,
	}
}

// HomePage serves the main upload page
func (h *UploadHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Pass the full config struct instead of anonymous struct
	if err := h.template.Execute(w, h.config); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// UploadFile handles file upload
func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Set max memory for multipart parsing (32MB)
	r.Body = http.MaxBytesReader(w, r.Body, h.config.MaxFileSize)

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, false, "File too large or invalid form data", nil)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.sendJSONResponse(w, http.StatusBadRequest, false, "No file uploaded", nil)
		return
	}
	defer file.Close()

	// Validate file size
	if handler.Size > h.config.MaxFileSize {
		h.sendJSONResponse(w, http.StatusBadRequest, false, "File exceeds maximum size limit", nil)
		return
	}

	// Sanitize filename
	filename := h.sanitizeFilename(handler.Filename)
	if filename == "" {
		h.sendJSONResponse(w, http.StatusBadRequest, false, "Invalid filename", nil)
		return
	}

	// Create unique filename if file already exists
	finalPath := h.getUniqueFilePath(filename)

	// Create the file
	dst, err := os.Create(finalPath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		h.sendJSONResponse(w, http.StatusInternalServerError, false, "Failed to save file", nil)
		return
	}
	defer dst.Close()

	// Copy file contents
	_, err = io.Copy(dst, file)
	if err != nil {
		log.Printf("Failed to copy file: %v", err)
		os.Remove(finalPath) // Clean up on error
		h.sendJSONResponse(w, http.StatusInternalServerError, false, "Failed to save file", nil)
		return
	}

	// Get file info for logging
	if stat, err := dst.Stat(); err != nil {
		log.Printf("Failed to get file stats: %v", err)
	} else {
		log.Printf("File saved successfully: %s (%d bytes on disk)", stat.Name(), stat.Size())
	}

	fileInfo := &FileInfo{
		Name:       filepath.Base(finalPath),
		Size:       handler.Size,
		UploadTime: time.Now(),
	}

	log.Printf("File uploaded successfully: %s (%d bytes)", fileInfo.Name, fileInfo.Size)
	h.sendJSONResponse(w, http.StatusOK, true, "File uploaded successfully", fileInfo)
}

// ListFiles returns list of uploaded files
func (h *UploadHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(h.config.DataDir)
	if err != nil {
		log.Printf("Failed to read data directory: %v", err)
		h.sendJSONResponse(w, http.StatusInternalServerError, false, "Failed to list files", nil)
		return
	}

	var fileList []FileInfo
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}
			fileList = append(fileList, FileInfo{
				Name:       file.Name(),
				Size:       info.Size(),
				UploadTime: info.ModTime(),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)
}

// sanitizeFilename removes unsafe characters from filename
func (h *UploadHandler) sanitizeFilename(filename string) string {
	// Remove path separators and other unsafe characters
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.TrimSpace(filename)

	if filename == "" || filename == "." {
		return ""
	}

	return filename
}

// getUniqueFilePath generates a unique file path to avoid conflicts
func (h *UploadHandler) getUniqueFilePath(filename string) string {
	basePath := filepath.Join(h.config.DataDir, filename)

	// If file doesn't exist, use original name
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return basePath
	}

	// Generate unique name with timestamp
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Format("20060102_150405")

	return filepath.Join(h.config.DataDir, fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext))
}

// sendJSONResponse sends a JSON response
func (h *UploadHandler) sendJSONResponse(w http.ResponseWriter, statusCode int, success bool, message string, file *FileInfo) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := UploadResponse{
		Success: success,
		Message: message,
		File:    file,
	}

	json.NewEncoder(w).Encode(response)
}
