package services

import (
	"embed"
	"errors"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/snowlynxsoftware/parallax-game/server/util"
)

//go:embed static/css/* static/js/* static/images/*
var staticFS embed.FS

type IStaticService interface {
	ServeStaticFile(w http.ResponseWriter, r *http.Request, filePath string) error
}

type StaticService struct {
}

func NewStaticService() IStaticService {
	return &StaticService{}
}

func (ss *StaticService) ServeStaticFile(w http.ResponseWriter, r *http.Request, filePath string) error {
	// Clean the file path to prevent directory traversal
	cleanPath := path.Clean(filePath)
	if strings.Contains(cleanPath, "..") {
		util.LogError(errors.New("Invalid file path: " + filePath))
		http.Error(w, "Forbidden", http.StatusForbidden)
		return errors.New("invalid file path")
	}

	// Remove leading slash and add static prefix
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	if !strings.HasPrefix(cleanPath, "static/") {
		cleanPath = "static/" + cleanPath
	}

	// Read the file from embedded filesystem
	content, err := staticFS.ReadFile(cleanPath)
	if err != nil {
		util.LogError(errors.New("Static file not found: " + cleanPath + " - " + err.Error()))
		http.Error(w, "Not Found", http.StatusNotFound)
		return err
	}

	// Determine content type based on file extension
	ext := filepath.Ext(cleanPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// Default content types for common web assets
		switch ext {
		case ".css":
			contentType = "text/css"
		case ".js":
			contentType = "application/javascript"
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".gif":
			contentType = "image/gif"
		case ".svg":
			contentType = "image/svg+xml"
		case ".ico":
			contentType = "image/x-icon"
		case ".woff":
			contentType = "font/woff"
		case ".woff2":
			contentType = "font/woff2"
		case ".ttf":
			contentType = "font/ttf"
		case ".eot":
			contentType = "application/vnd.ms-fontobject"
		default:
			contentType = "application/octet-stream"
		}
	}

	// Set headers
	w.Header().Set("Content-Type", contentType)

	// Cache static assets for 1 hour
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// Add security headers for certain file types
	if ext == ".js" || ext == ".css" {
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// Write the content
	_, err = w.Write(content)
	if err != nil {
		util.LogError(errors.New("Failed to write static file content: " + err.Error()))
		return err
	}

	util.LogDebug("Served static file: " + cleanPath)
	return nil
}
