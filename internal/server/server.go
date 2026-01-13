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
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/journal"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
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
nav { border-bottom: 1px solid #ddd; margin-bottom: 20px; }
article {
    background: white;
    padding: 20px;
    border-radius: 5px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.1);
    margin-bottom: 20px;
}
article h4 { margin-top: 0; }
article pre, article code { background: #f4f4f4; }
article pre { padding: 15px; border-radius: 5px; overflow-x: auto; }
article code { padding: 2px 6px; border-radius: 3px; }
article pre code { background: none; padding: 0; }
article blockquote { border-left: 4px solid #ddd; margin: 0; padding-left: 20px; color: #666; }
`

// Server represents the journal preview server
type Server struct {
	cfg        *config.Config
	journal    *journal.Journal
	baseDir    string
	css        string
	liveReload bool

	mu      sync.RWMutex
	entries journal.Entries
	tmpl    *template.Template
	md      goldmark.Markdown

	// SSE clients for live reload
	sseClients   map[chan struct{}]struct{}
	sseClientsMu sync.Mutex
}

// New creates a new Server instance
func New(cfg *config.Config, jnl *journal.Journal, baseDir string, liveReload bool) (*Server, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	// Load CSS
	css, err := loadCSS(cfg.Build.CSS)
	if err != nil {
		return nil, fmt.Errorf("loading css: %w", err)
	}

	// Configure goldmark
	md := goldmark.New()
	if cfg.Build.GetHardWraps() {
		md = goldmark.New(
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
			),
		)
	}

	return &Server{
		cfg:        cfg,
		journal:    jnl,
		baseDir:    baseDir,
		css:        css,
		liveReload: liveReload,
		tmpl:       tmpl,
		md:         md,
		sseClients: make(map[chan struct{}]struct{}),
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
	if s.liveReload {
		mux.HandleFunc("/events", s.handleSSE)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.Serve.Port),
		Handler: mux,
	}

	// Handle graceful shutdown
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	fmt.Printf("Starting server at http://localhost:%d\n", s.cfg.Serve.Port)
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
	switch s.cfg.Build.Sort {
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

	shift := s.cfg.Build.GetHeadingShift()
	if shift > 0 {
		return shiftHeadings(buf.String(), shift), nil
	}
	return buf.String(), nil
}

// shiftHeadings shifts HTML heading levels by the specified amount
// H1 becomes H1+shift, H2 becomes H2+shift, etc.
// Headings are clamped to H6 maximum
func shiftHeadings(html string, shift int) string {
	// Process from h6 to h1 to avoid double replacement
	for level := 6; level >= 1; level-- {
		newLevel := level + shift
		if newLevel > 6 {
			newLevel = 6
		}

		// Use regex to handle attributes in opening tags
		re := regexp.MustCompile(fmt.Sprintf(`<h%d(\s|>)`, level))
		html = re.ReplaceAllString(html, fmt.Sprintf("<h%d$1", newLevel))
		html = strings.ReplaceAll(html, fmt.Sprintf("</h%d>", level), fmt.Sprintf("</h%d>", newLevel))
	}
	return html
}

// watchFiles watches for file changes and reloads entries
func (s *Server) watchFiles(ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Watch base directory and all subdirectories
	err = filepath.Walk(s.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := watcher.Add(path); err != nil {
				fmt.Printf("Error watching directory %s: %v\n", path, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
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
					s.notifyClients()
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

// handleSSE handles Server-Sent Events for live reload
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a channel for this client
	clientChan := make(chan struct{}, 1)

	// Register client
	s.sseClientsMu.Lock()
	s.sseClients[clientChan] = struct{}{}
	s.sseClientsMu.Unlock()

	// Unregister client on disconnect
	defer func() {
		s.sseClientsMu.Lock()
		delete(s.sseClients, clientChan)
		s.sseClientsMu.Unlock()
		close(clientChan)
	}()

	// Send initial connection message
	fmt.Fprintf(w, "data: connected\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Wait for reload signals or client disconnect
	for {
		select {
		case <-clientChan:
			fmt.Fprintf(w, "data: reload\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

// notifyClients sends reload signal to all connected SSE clients
func (s *Server) notifyClients() {
	s.sseClientsMu.Lock()
	defer s.sseClientsMu.Unlock()

	for clientChan := range s.sseClients {
		select {
		case clientChan <- struct{}{}:
		default:
			// Channel is full, skip
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

	templateEntries, yearNavs := convertToTemplateEntries(entries)

	data := IndexData{
		Title:      s.cfg.Build.Title,
		Entries:    templateEntries,
		YearNavs:   yearNavs,
		CSS:        template.CSS(s.css),
		LiveReload: s.liveReload,
	}

	if err := s.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// convertToTemplateEntries converts journal entries to template entries with year/month markers
func convertToTemplateEntries(entries journal.Entries) ([]TemplateEntry, []YearNav) {
	templateEntries := make([]TemplateEntry, len(entries))
	yearNavs := []YearNav{}
	lastYear := ""
	lastMonth := ""

	for i, e := range entries {
		year := e.Date.Format("2006")
		month := e.Date.Format("01")
		yearMonth := year + "-" + month

		showYear := year != lastYear
		showMonth := yearMonth != lastMonth

		if showYear {
			yearNavs = append(yearNavs, YearNav{Year: year, Months: []string{month}})
			lastYear = year
		} else if showMonth {
			yearNavs[len(yearNavs)-1].Months = append(yearNavs[len(yearNavs)-1].Months, month)
		}
		lastMonth = yearMonth

		templateEntries[i] = TemplateEntry{
			Date:       e.Date,
			Content:    template.HTML(e.Content),
			ShowYear:   showYear,
			YearLabel:  year,
			ShowMonth:  showMonth,
			MonthLabel: yearMonth,
		}
	}

	return templateEntries, yearNavs
}

// TemplateEntry represents an entry for template rendering
type TemplateEntry struct {
	Date       time.Time
	Content    template.HTML
	ShowYear   bool
	YearLabel  string
	ShowMonth  bool
	MonthLabel string
}

// YearNav represents navigation for a year
type YearNav struct {
	Year   string
	Months []string
}

// IndexData represents data for the index template
type IndexData struct {
	Title      string
	Entries    []TemplateEntry
	YearNavs   []YearNav
	CSS        template.CSS
	LiveReload bool
}

// Builder generates static HTML files
type Builder struct {
	cfg     *config.Config
	journal *journal.Journal
	baseDir string
	css     string
	tmpl    *template.Template
	md      goldmark.Markdown
}

// NewBuilder creates a new Builder instance
func NewBuilder(cfg *config.Config, jnl *journal.Journal, baseDir string) (*Builder, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parsing templates: %w", err)
	}

	css, err := loadCSS(cfg.Build.CSS)
	if err != nil {
		return nil, fmt.Errorf("loading css: %w", err)
	}

	// Configure goldmark
	md := goldmark.New()
	if cfg.Build.GetHardWraps() {
		md = goldmark.New(
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
			),
		)
	}

	return &Builder{
		cfg:     cfg,
		journal: jnl,
		baseDir: baseDir,
		css:     css,
		tmpl:    tmpl,
		md:      md,
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
	switch b.cfg.Build.Sort {
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

	// Convert to template entries
	templateEntries, yearNavs := convertToTemplateEntries(entries)

	// Generate index.html
	indexData := IndexData{
		Title:    b.cfg.Build.Title,
		Entries:  templateEntries,
		YearNavs: yearNavs,
		CSS:      template.CSS(b.css),
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

	shift := b.cfg.Build.GetHeadingShift()
	if shift > 0 {
		return shiftHeadings(buf.String(), shift), nil
	}
	return buf.String(), nil
}
