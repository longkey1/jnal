package server

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/journal"
	"github.com/yuin/goldmark"
)

//go:embed templates/*.html
var templatesFS embed.FS

// DefaultCSS is the default stylesheet
const DefaultCSS = `
* { box-sizing: border-box; }
body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    background-color: #fafafa;
    color: #333;
}
h1 { border-bottom: 2px solid #333; padding-bottom: 10px; }
a { color: #007acc; text-decoration: none; }
a:hover { text-decoration: underline; }
ul { list-style: none; padding: 0; }
li { padding: 8px 0; border-bottom: 1px solid #eee; }
li:last-child { border-bottom: none; }
header { margin-bottom: 20px; padding-bottom: 10px; border-bottom: 1px solid #ddd; }
.content {
    background: white;
    padding: 20px;
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}
.content h1:first-child { margin-top: 0; }
.content pre, .content code { background: #f4f4f4; }
.content pre { padding: 15px; border-radius: 5px; overflow-x: auto; }
.content code { padding: 2px 6px; border-radius: 3px; }
.content pre code { background: none; padding: 0; }
.content blockquote { border-left: 4px solid #ddd; margin: 0; padding-left: 20px; color: #666; }
.meta { color: #666; font-size: 0.9em; margin-bottom: 20px; }
.date { font-family: monospace; color: #666; margin-right: 10px; }
`

// Server represents the journal preview server
type Server struct {
	cfg     *config.ServeConfig
	journal *journal.Journal
	baseDir string
	css     string

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

	// Load CSS
	css, err := loadCSS(cfg.CSS)
	if err != nil {
		return nil, fmt.Errorf("loading css: %w", err)
	}

	return &Server{
		cfg:     cfg,
		journal: jnl,
		baseDir: baseDir,
		css:     css,
		tmpl:    tmpl,
		md:      goldmark.New(),
	}, nil
}

// loadCSS loads CSS from URL or returns the string as-is
func loadCSS(css string) (string, error) {
	if css == "" {
		return DefaultCSS, nil
	}

	// Check if it's a URL
	if strings.HasPrefix(css, "http://") || strings.HasPrefix(css, "https://") {
		resp, err := http.Get(css)
		if err != nil {
			return "", fmt.Errorf("fetching CSS from %s: %w", css, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("fetching CSS from %s: status %d", css, resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("reading CSS response: %w", err)
		}

		return string(body), nil
	}

	// Return as-is (inline CSS)
	return css, nil
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
		CSS:     template.CSS(s.css),
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
		CSS:     template.CSS(s.css),
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
	CSS     template.CSS
}

// EntryData represents data for the entry template
type EntryData struct {
	Title   string
	Entry   journal.Entry
	Content template.HTML
	CSS     template.CSS
}

// Builder generates static HTML files
type Builder struct {
	cfg     *config.ServeConfig
	journal *journal.Journal
	baseDir string
	css     string
	tmpl    *template.Template
	md      goldmark.Markdown
}

// NewBuilder creates a new Builder instance
func NewBuilder(cfg *config.ServeConfig, jnl *journal.Journal, baseDir string) (*Builder, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	css, err := loadCSS(cfg.CSS)
	if err != nil {
		return nil, fmt.Errorf("loading css: %w", err)
	}

	return &Builder{
		cfg:     cfg,
		journal: jnl,
		baseDir: baseDir,
		css:     css,
		tmpl:    tmpl,
		md:      goldmark.New(),
	}, nil
}

// Build generates static HTML files to the output directory
func (b *Builder) Build(outputDir string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Load entries
	entries, err := b.journal.ListEntries()
	if err != nil {
		return fmt.Errorf("listing entries: %w", err)
	}

	// Sort entries
	switch b.cfg.Sort {
	case config.SortAsc:
		entries.SortByDateAsc()
	default:
		entries.SortByDateDesc()
	}

	// Load content for each entry
	for i := range entries {
		content, err := b.loadEntryContent(entries[i].Path)
		if err != nil {
			continue
		}
		entries[i].Content = content
	}

	// Generate index.html
	indexData := IndexData{
		Title:   "Journal",
		Entries: entries,
		Sort:    b.cfg.Sort,
		Updated: time.Now(),
		CSS:     template.CSS(b.css),
	}

	indexPath := filepath.Join(outputDir, "index.html")
	indexFile, err := os.Create(indexPath)
	if err != nil {
		return fmt.Errorf("creating index.html: %w", err)
	}
	defer indexFile.Close()

	if err := b.tmpl.ExecuteTemplate(indexFile, "index.html", indexData); err != nil {
		return fmt.Errorf("executing index template: %w", err)
	}

	// Generate entry pages
	entryDir := filepath.Join(outputDir, "entry")
	if err := os.MkdirAll(entryDir, 0755); err != nil {
		return fmt.Errorf("creating entry directory: %w", err)
	}

	for _, entry := range entries {
		entryData := EntryData{
			Title:   entry.Date.Format("2006-01-02"),
			Entry:   entry,
			Content: template.HTML(entry.Content),
			CSS:     template.CSS(b.css),
		}

		entryPath := filepath.Join(entryDir, entry.Date.Format("2006-01-02")+".html")
		entryFile, err := os.Create(entryPath)
		if err != nil {
			return fmt.Errorf("creating entry file: %w", err)
		}

		if err := b.tmpl.ExecuteTemplate(entryFile, "entry.html", entryData); err != nil {
			entryFile.Close()
			return fmt.Errorf("executing entry template: %w", err)
		}
		entryFile.Close()
	}

	return nil
}

// loadEntryContent loads and converts markdown content to HTML
func (b *Builder) loadEntryContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := b.md.Convert(data, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
