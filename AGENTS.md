# AGENTS.md

This file provides guidance to AI coding agents (Claude Code, etc.) when working with code in this repository.

## What This Project Does

**jnal** is a CLI tool for managing Markdown-based journal entries with static HTML generation and a local preview server.

## Common Commands

```bash
make build          # Build binary to ./bin/
make test           # Run tests (go test ./...)
make fmt            # Format code (go fmt ./...)
make vet            # Vet code (go vet ./...)
make tidy           # Tidy dependencies (go mod tidy)
make clean          # Clean build artifacts
```

Run a single test package:
```bash
go test ./internal/config/...
go test ./internal/util/...
```

Release management:
```bash
make release type=patch dryrun=true    # Dry run
make release type=patch dryrun=false   # v1.2.3 → v1.2.4
make re-release tag=v1.2.3 dryrun=false
```

## Architecture

### Package Structure

```
main.go              # Entrypoint: calls cmd.NewRootCommand().Execute()
cmd/                 # Cobra CLI commands (new, build, serve, path, init, version)
internal/
  jnal/              # Core domain: App, Journal, Entry structs
  config/            # Config struct, Viper-based TOML loader, validation
  server/            # HTTP server + static builder, embedded HTML template
  util/              # Date parsing/formatting utilities
  version/           # Version info injected via ldflags at build time
```

### Key Abstractions

**`internal/jnal`** — Core application state:
- `App`: Top-level container holding `*config.Config` and `*Journal`. Created in `root.go`'s `PersistentPreRunE` hook and passed to each subcommand.
- `Journal`: Handles entry CRUD — `CreateEntry`, `GetEntryPath`, `ListEntries`, `EntryExists`. Path generation is controlled by `config.Common.PathFormat` (a Go time format string, e.g. `"2006/01/2006-01-02.md"`).
- `Entry` / `Entries`: Value type with `Path`, `Date`, `Content` (HTML). Entries can be sorted by `SortByDateAsc()` / `SortByDateDesc()`.

**`internal/config`** — Configuration:
- Default config path: `$HOME/.config/jnal/config.toml` (overridable via `JNAL_CONFIG` env or `--config` flag).
- `Config` has four sections: `Common`, `New`, `Build`, `Serve`.
- Optional fields use pointer types (`*int`, `*bool`) — check for nil before dereferencing.
- Each section has `Validate()` and `SetDefaults()` methods called during load.

**`internal/server`** — Rendering:
- `Server`: HTTP server with optional live-reload via `fsnotify` + Server-Sent Events. Caches rendered HTML in `entries []jnal.Entries` protected by `sync.RWMutex`.
- `Builder`: Same rendering logic as `Server` but writes a single `public/index.html` instead of serving.
- HTML template is embedded with `//go:embed templates/index.html`.
- Markdown rendering uses Goldmark; `shiftHeadings()` and `addTargetBlankToLinks()` post-process the HTML output.

### Command → Code Flow

```
jnal new [--date]
  root PersistentPreRunE → jnal.NewApp(cfgFile) → config.Load() → Journal
  new RunE → util.Parse(date) → journal.CreateEntry(date) → template execution → file write

jnal build [--output]
  → server.NewBuilder(cfg, journal) → builder.Build(outputDir) → index.html

jnal serve [--port] [--sort] [--live-reload]
  → server.NewServer(cfg, journal) → server.Start(ctx) → HTTP + optional fsnotify watcher
```

### Configuration Sections

| Section | Key Fields |
|---------|-----------|
| `common` | `base_directory` (required), `date_format`, `path_format` |
| `new` | `file_template` (supports `{{ .Date }}` and `{{ .Env.<NAME> }}`) |
| `build` | `title`, `sort` (`"asc"`/`"desc"`), `css`, `heading_shift`, `hard_wraps`, `linkify`, `link_target_blank` |
| `serve` | `port` (default: 8080) |
