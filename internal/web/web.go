package web

import (
	"embed"
	"errors"
	"github.com/xuanyiying/entimport/internal/entimport"
	"github.com/xuanyiying/entimport/internal/mux"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed templates/*.html
var templatesFS embed.FS

// Server represents the web server for entimport UI
type Server struct {
	templates *template.Template
	config    *Config
}

// Config holds server configuration
type Config struct {
	TemplatesPath string // Path to template directory
}

// NewServer creates a new web server instance
func NewServer(cfg *Config) (*Server, error) {
	s := &Server{
		config: cfg,
	}

	// Try loading templates from embedded FS first
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err == nil {
		s.templates = tmpl
		return s, nil
	}

	// If custom template path is provided, try loading from there
	if cfg != nil && cfg.TemplatesPath != "" {
		tmpl, err := template.ParseGlob(filepath.Join(cfg.TemplatesPath, "*.html"))
		if err != nil {
			return nil, err
		}
		s.templates = tmpl
		return s, nil
	}

	return nil, err
}

// Start starts the web server
func (s *Server) Start(addr string) error {
	if s.templates == nil {
		return errors.New("templates not initialized")
	}

	http.HandleFunc("/", s.handleHome)
	http.HandleFunc("/import", s.handleImport)
	return http.ListenAndServe(addr, nil)
}

// handleHome renders the home page
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	s.templates.ExecuteTemplate(w, "index.html", nil)
}

// ImportRequest represents the form data for import
type ImportRequest struct {
	DSN           string   `json:"dsn"`
	Tables        []string `json:"tables"`
	ExcludeTables []string `json:"exclude_tables"`
	SchemaPath    string   `json:"schema_path"`
}

// handleImport handles the import request
func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ImportRequest
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse form data
	req.DSN = r.FormValue("dsn")
	req.SchemaPath = r.FormValue("schema_path")
	if tables := r.FormValue("tables"); tables != "" {
		req.Tables = strings.Split(tables, ",")
	}
	if excludeTables := r.FormValue("exclude_tables"); excludeTables != "" {
		req.ExcludeTables = strings.Split(excludeTables, ",")
	}

	// Create import driver
	drv, err := mux.Default.OpenImport(req.DSN)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create importer
	importer, err := entimport.NewImport(
		entimport.WithDriver(drv),
		entimport.WithSchemaPath(req.SchemaPath),
		entimport.WithTables(req.Tables),
		entimport.WithExcludedTables(req.ExcludeTables),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Run import
	mutations, err := importer.SchemaMutations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := entimport.WriteSchema(mutations, entimport.WithSchemaPath(req.SchemaPath)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Import successful"))
}
