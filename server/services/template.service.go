package services

import (
	"embed"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/snowlynxsoftware/parallax-game/server/util"
)

//go:embed templates/*
var templatesFS embed.FS

type ITemplateService interface {
	RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error
	GetTemplate(templateName string) (*template.Template, error)
}

type TemplateService struct {
	templates map[string]*template.Template
}

type PageData struct {
	Title       string
	Description string
	Data        interface{}
}

func NewTemplateService() ITemplateService {
	service := &TemplateService{
		templates: make(map[string]*template.Template),
	}
	service.loadTemplates()
	return service
}

func (ts *TemplateService) loadTemplates() {
	// Define template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"multiply": func(a, b int) int {
			return a * b
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"formatDate": func(dateStr string) string {
			return formatDate(dateStr)
		},
	}

	// Load base layout and components
	baseTemplate := template.New("base").Funcs(funcMap)

	// Parse base layout
	baseContent, err := templatesFS.ReadFile("templates/layouts/base.html")
	if err != nil {
		util.LogError(errors.New("Failed to read base template: " + err.Error()))
		return
	}

	baseTemplate, err = baseTemplate.Parse(string(baseContent))
	if err != nil {
		util.LogError(errors.New("Failed to parse base template: " + err.Error()))
		return
	}

	// Parse components
	componentFiles := []string{
		"templates/layouts/components/head.html",
		"templates/layouts/components/footer.html",
		"templates/layouts/components/navbar.html",
	}

	for _, file := range componentFiles {
		content, err := templatesFS.ReadFile(file)
		if err != nil {
			util.LogError(errors.New("Failed to read component template " + file + ": " + err.Error()))
			continue
		}

		_, err = baseTemplate.Parse(string(content))
		if err != nil {
			util.LogError(errors.New("Failed to parse component template " + file + ": " + err.Error()))
			continue
		}
	}

	// Load page templates
	pageFiles := []string{
		"templates/pages/welcome.html",
		"templates/pages/register.html",
		"templates/pages/login.html",
		"templates/pages/dashboard.html",
		"templates/pages/teams.html",
		"templates/pages/expeditions.html",
		"templates/pages/inventory.html",
		"templates/pages/account.html",
		"templates/pages/reset-password.html",
		"templates/pages/terms.html",
		"templates/pages/privacy.html",
	}

	for _, file := range pageFiles {
		content, err := templatesFS.ReadFile(file)
		if err != nil {
			util.LogError(errors.New("Failed to read page template " + file + ": " + err.Error()))
			continue
		}

		// Clone the base template for each page
		pageTemplate, err := baseTemplate.Clone()
		if err != nil {
			util.LogError(errors.New("Failed to clone base template for " + file + ": " + err.Error()))
			continue
		}

		_, err = pageTemplate.Parse(string(content))
		if err != nil {
			util.LogError(errors.New("Failed to parse page template " + file + ": " + err.Error()))
			continue
		}

		// Extract template name from file path
		templateName := strings.TrimSuffix(filepath.Base(file), ".html")
		ts.templates[templateName] = pageTemplate

		util.LogDebug("Loaded template: " + templateName)
	}
}

func (ts *TemplateService) GetTemplate(templateName string) (*template.Template, error) {
	tmpl, exists := ts.templates[templateName]
	if !exists {
		util.LogError(errors.New("Template not found: " + templateName))
		return nil, errors.New("template not found")
	}
	return tmpl, nil
}

func (ts *TemplateService) RenderTemplate(w http.ResponseWriter, templateName string, data interface{}) error {
	tmpl, err := ts.GetTemplate(templateName)
	if err != nil {
		util.LogError(errors.New("Failed to get template " + templateName + ": " + err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	if tmpl == nil {
		util.LogError(errors.New("Template not found: " + templateName))
		http.Error(w, "Template Not Found", http.StatusNotFound)
		return errors.New("template not found")
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		util.LogError(errors.New("Failed to execute template " + templateName + ": " + err.Error()))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}

	return nil
}

// formatDate formats an RFC3339 date string to a more readable format
func formatDate(dateStr string) string {
	if dateStr == "" {
		return "Never"
	}

	// Parse the RFC3339 date string
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return dateStr // Return original if parsing fails
	}

	// Format as "Jan 02, 2006 15:04"
	return t.Format("Jan 02, 2006 15:04")
}
