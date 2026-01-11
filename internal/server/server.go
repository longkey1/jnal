package server

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/journal"
	"github.com/yuin/goldmark"
)

//go:embed templates/*.html
var templatesFS embed.FS

// Server represents the journal preview server
type Server struct {
	cfg     *config.ServeConfig
	journal *journal.Journal
	baseDir string

	mu      sync.RWMutex
	entries journal.Entries
	tmpl    *template.Template
	md      goldmark.Markdown
}

// New creates a new Server instance
func New(cfg *config.ServeConfig, jnl *journal.Journal, baseDir string) (*Server, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	return &Server{
		cfg:     cfg,
		journal: jnl,
		baseDir: baseDir,
		tmpl:    tmpl,
		md:      goldmark.New(),
	}, nil
}

// Start starts the server
func (s *Server) Start(ctx context.Context) error {
	// Load initial entries
	if err := s.reloadEntries(); err != nil {
		return fmt.Errorf("loading entries: %w", err)
	}

	// Start file watcher
	go s.watchFiles(ctx)

	// Setup HTTP handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/entry/", s.handleEntry)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Port),
		Handler: mux,
	}

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("Starting server at http://localhost:%d\n", s.cfg.Port)
	fmt.Println("Press Ctrl+C to stop")

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// reloadEntries reloads journal entries from disk
func (s *Server) reloadEntries() error {
	entries, err := s.journal.ListEntries()
	if err != nil {
		return err
	}

	// Sort entries based on configuration
	switch s.cfg.Sort {
	case config.SortAsc:
		entries.SortByDateAsc()
	default:
		entries.SortByDateDesc()
	}

	// Load content for each entry
	for i := range entries {
		content, err := s.loadEntryContent(entries[i].Path)
		if err != nil {
			continue
		}
		entries[i].Content = content
	}

	s.mu.Lock()
	s.entries = entries
	s.mu.Unlock()

	return nil
}

// loadEntryContent loads and converts markdown content to HTML
func (s *Server) loadEntryContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := s.md.Convert(data, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// watchFiles watches for file changes and reloads entries
func (s *Server) watchFiles(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	if err := watcher.Add(s.baseDir); err != nil {
		fmt.Printf("Error watching directory: %v\n", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) != 0 {
				if filepath.Ext(event.Name) == ".md" {
					fmt.Printf("File changed: %s, reloading...\n", filepath.Base(event.Name))
					if err := s.reloadEntries(); err != nil {
						fmt.Printf("Error reloading entries: %v\n", err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// handleIndex handles the index page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	s.mu.RLock()
	entries := s.entries
	s.mu.RUnlock()

	data := IndexData{
		Title:   "Journal",
		Entries: entries,
		Sort:    s.cfg.Sort,
		Updated: time.Now(),
	}

	if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleEntry handles individual entry pages
func (s *Server) handleEntry(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Path[len("/entry/"):]
	if dateStr == "" {
		http.NotFound(w, r)
		return
	}

	// Remove .html suffix if present
	if len(dateStr) > 5 && dateStr[len(dateStr)-5:] == ".html" {
		dateStr = dateStr[:len(dateStr)-5]
	}

	s.mu.RLock()
	var entry *journal.Entry
	for i := range s.entries {
		if s.entries[i].Date.Format("2006-01-02") == dateStr {
			entry = &s.entries[i]
			break
		}
	}
	s.mu.RUnlock()

	if entry == nil {
		http.NotFound(w, r)
		return
	}

	data := EntryData{
		Title:   entry.Date.Format("2006-01-02"),
		Entry:   *entry,
		Content: template.HTML(entry.Content),
	}

	if err := s.tmpl.ExecuteTemplate(w, "entry.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// IndexData represents data for the index template
type IndexData struct {
	Title   string
	Entries journal.Entries
	Sort    string
	Updated time.Time
}

// EntryData represents data for the entry template
type EntryData struct {
	Title   string
	Entry   journal.Entry
	Content template.HTML
}
