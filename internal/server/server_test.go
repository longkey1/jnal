package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/longkey1/jnal/internal/config"
	"github.com/longkey1/jnal/internal/jnal"
)

func TestShiftHeadings(t *testing.T) {
	tests := []struct {
		name  string
		html  string
		shift int
		want  string
	}{
		{
			name:  "shift h1 to h5",
			html:  "<h1>Title</h1>",
			shift: 4,
			want:  "<h5>Title</h5>",
		},
		{
			name:  "shift multiple heading levels",
			html:  "<h1>A</h1><h2>B</h2>",
			shift: 1,
			want:  "<h2>A</h2><h3>B</h3>",
		},
		{
			name:  "clamp at h6",
			html:  "<h3>Deep</h3><h6>Deepest</h6>",
			shift: 4,
			want:  "<h6>Deep</h6><h6>Deepest</h6>",
		},
		{
			name:  "heading with attributes",
			html:  `<h1 id="title" class="main">Title</h1>`,
			shift: 2,
			want:  `<h3 id="title" class="main">Title</h3>`,
		},
		{
			name:  "no double replacement of shifted headings",
			html:  "<h1>A</h1><h5>B</h5>",
			shift: 4,
			want:  "<h5>A</h5><h6>B</h6>",
		},
		{
			name:  "zero shift leaves html unchanged",
			html:  "<h2>Title</h2>",
			shift: 0,
			want:  "<h2>Title</h2>",
		},
		{
			name:  "non-heading content untouched",
			html:  "<p>h1 is not a tag here</p>",
			shift: 3,
			want:  "<p>h1 is not a tag here</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shiftHeadings(tt.html, tt.shift)
			if got != tt.want {
				t.Errorf("shiftHeadings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTargetBlankToLinks(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		{
			name: "single link",
			html: `<a href="https://example.com">example</a>`,
			want: `<a target="_blank" rel="noopener noreferrer" href="https://example.com">example</a>`,
		},
		{
			name: "multiple links",
			html: `<a href="/a">a</a> and <a href="/b">b</a>`,
			want: `<a target="_blank" rel="noopener noreferrer" href="/a">a</a> and <a target="_blank" rel="noopener noreferrer" href="/b">b</a>`,
		},
		{
			name: "no links",
			html: "<p>plain text</p>",
			want: "<p>plain text</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addTargetBlankToLinks(tt.html)
			if got != tt.want {
				t.Errorf("addTargetBlankToLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadCSS(t *testing.T) {
	tests := []struct {
		name string
		css  string
		want string
	}{
		{
			name: "empty returns default css",
			css:  "",
			want: DefaultCSS,
		},
		{
			name: "inline css returned as-is",
			css:  "body { color: red; }",
			want: "body { color: red; }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadCSS(tt.css)
			if err != nil {
				t.Fatalf("loadCSS() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("loadCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func entryOn(y int, m time.Month, d int) jnal.Entry {
	return jnal.Entry{Date: time.Date(y, m, d, 0, 0, 0, 0, time.UTC)}
}

func TestConvertToTemplateEntries(t *testing.T) {
	tests := []struct {
		name          string
		entries       jnal.Entries
		wantShowYear  []bool
		wantShowMonth []bool
		wantYearNavs  []YearNav
	}{
		{
			name:          "empty entries",
			entries:       jnal.Entries{},
			wantShowYear:  []bool{},
			wantShowMonth: []bool{},
			wantYearNavs:  []YearNav{},
		},
		{
			name: "same month entries share markers",
			entries: jnal.Entries{
				entryOn(2024, 1, 15),
				entryOn(2024, 1, 10),
			},
			wantShowYear:  []bool{true, false},
			wantShowMonth: []bool{true, false},
			wantYearNavs: []YearNav{
				{Year: "2024", Months: []string{"01"}},
			},
		},
		{
			name: "month change within same year",
			entries: jnal.Entries{
				entryOn(2024, 2, 1),
				entryOn(2024, 1, 31),
			},
			wantShowYear:  []bool{true, false},
			wantShowMonth: []bool{true, true},
			wantYearNavs: []YearNav{
				{Year: "2024", Months: []string{"02", "01"}},
			},
		},
		{
			name: "year change starts new nav",
			entries: jnal.Entries{
				entryOn(2024, 1, 1),
				entryOn(2023, 12, 31),
				entryOn(2023, 11, 1),
			},
			wantShowYear:  []bool{true, true, false},
			wantShowMonth: []bool{true, true, true},
			wantYearNavs: []YearNav{
				{Year: "2024", Months: []string{"01"}},
				{Year: "2023", Months: []string{"12", "11"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			templateEntries, yearNavs := convertToTemplateEntries(tt.entries)

			if len(templateEntries) != len(tt.entries) {
				t.Fatalf("convertToTemplateEntries() returned %d entries, want %d", len(templateEntries), len(tt.entries))
			}
			for i, te := range templateEntries {
				if te.ShowYear != tt.wantShowYear[i] {
					t.Errorf("entry[%d].ShowYear = %v, want %v", i, te.ShowYear, tt.wantShowYear[i])
				}
				if te.ShowMonth != tt.wantShowMonth[i] {
					t.Errorf("entry[%d].ShowMonth = %v, want %v", i, te.ShowMonth, tt.wantShowMonth[i])
				}
				wantYearLabel := tt.entries[i].Date.Format("2006")
				if te.YearLabel != wantYearLabel {
					t.Errorf("entry[%d].YearLabel = %v, want %v", i, te.YearLabel, wantYearLabel)
				}
				wantMonthLabel := tt.entries[i].Date.Format("2006-01")
				if te.MonthLabel != wantMonthLabel {
					t.Errorf("entry[%d].MonthLabel = %v, want %v", i, te.MonthLabel, wantMonthLabel)
				}
			}

			if len(yearNavs) != len(tt.wantYearNavs) {
				t.Fatalf("convertToTemplateEntries() returned %d year navs, want %d", len(yearNavs), len(tt.wantYearNavs))
			}
			for i, nav := range yearNavs {
				want := tt.wantYearNavs[i]
				if nav.Year != want.Year {
					t.Errorf("yearNavs[%d].Year = %v, want %v", i, nav.Year, want.Year)
				}
				if len(nav.Months) != len(want.Months) {
					t.Fatalf("yearNavs[%d].Months = %v, want %v", i, nav.Months, want.Months)
				}
				for m := range want.Months {
					if nav.Months[m] != want.Months[m] {
						t.Errorf("yearNavs[%d].Months[%d] = %v, want %v", i, m, nav.Months[m], want.Months[m])
					}
				}
			}
		})
	}
}

func newBuilderTestConfig(baseDir string) *config.Config {
	cfg := &config.Config{
		Common: config.CommonConfig{
			BaseDirectory: baseDir,
		},
	}
	cfg.SetDefaults()
	return cfg
}

func TestBuilder_Build(t *testing.T) {
	baseDir := t.TempDir()
	outputDir := t.TempDir()

	markdown := "# Heading\n\nSee [example](https://example.com) for details.\n"
	if err := os.WriteFile(filepath.Join(baseDir, "2024-01-15.md"), []byte(markdown), 0644); err != nil {
		t.Fatalf("writing test entry: %v", err)
	}

	cfg := newBuilderTestConfig(baseDir)
	cfg.Build.Title = "Test Journal"
	jnl := jnal.NewJournal(cfg)

	builder, err := NewBuilder(cfg, jnl, baseDir)
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	if err := builder.Build(outputDir); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatalf("reading generated index.html: %v", err)
	}
	html := string(data)

	if !strings.Contains(html, "Test Journal") {
		t.Error("Build() output missing title")
	}
	// Default heading shift is 4, so H1 becomes H5
	if !strings.Contains(html, "<h5>Heading</h5>") {
		t.Error("Build() output missing shifted heading <h5>Heading</h5>")
	}
	// Default link_target_blank is true
	if !strings.Contains(html, `<a target="_blank" rel="noopener noreferrer" href="https://example.com">example</a>`) {
		t.Error("Build() output missing link with target blank")
	}
}

func TestBuilder_Build_HeadingShiftDisabled(t *testing.T) {
	baseDir := t.TempDir()
	outputDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(baseDir, "2024-01-15.md"), []byte("# Heading\n"), 0644); err != nil {
		t.Fatalf("writing test entry: %v", err)
	}

	cfg := newBuilderTestConfig(baseDir)
	noShift := 0
	cfg.Build.HeadingShift = &noShift
	jnl := jnal.NewJournal(cfg)

	builder, err := NewBuilder(cfg, jnl, baseDir)
	if err != nil {
		t.Fatalf("NewBuilder() error = %v", err)
	}

	if err := builder.Build(outputDir); err != nil {
		t.Fatalf("Build() error = %v", err)
	}

	data, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatalf("reading generated index.html: %v", err)
	}
	if !strings.Contains(string(data), "<h1>Heading</h1>") {
		t.Error("Build() output missing unshifted heading <h1>Heading</h1>")
	}
}
