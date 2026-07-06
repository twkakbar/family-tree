// Command migrate manages database migrations from the terminal:
//
//	go run ./cmd/migrate up          apply all pending migrations
//	go run ./cmd/migrate down [n]    roll back n migrations (default 1)
//	go run ./cmd/migrate version     print current version and dirty flag
//	go run ./cmd/migrate force <v>   mark version v as applied (repair dirty state)
//	go run ./cmd/migrate new <name>  scaffold numbered .up/.down SQL files
//
// Connection settings come from the environment / backend/.env, identical to
// the API server (internal/config).
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"

	"github.com/twkakbar/family-tree/backend/internal/config"
	"github.com/twkakbar/family-tree/backend/internal/repository"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: migrate <up|down [n]|version|force <v>|new <name>>")
	}

	// `new` only touches the filesystem; no DB or full config required.
	if args[0] == "new" {
		if len(args) < 2 {
			return errors.New("usage: migrate new <snake_case_name>")
		}
		return newMigration(args[1])
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := repository.Connect(cfg.DB, "warn")
	if err != nil {
		return err
	}
	m, err := repository.NewMigrator(db)
	if err != nil {
		return err
	}

	switch args[0] {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("already up to date")
				return nil
			}
			return err
		}
		fmt.Println("migrations applied")

	case "down":
		steps := 1
		if len(args) > 1 {
			if steps, err = strconv.Atoi(args[1]); err != nil || steps < 1 {
				return errors.New("down expects a positive step count")
			}
		}
		if err := m.Steps(-steps); err != nil {
			return err
		}
		fmt.Printf("rolled back %d migration(s)\n", steps)

	case "version":
		v, dirty, err := m.Version()
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("no migrations applied yet")
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("version=%d dirty=%v\n", v, dirty)

	case "force":
		if len(args) < 2 {
			return errors.New("usage: migrate force <version>")
		}
		v, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.New("force expects a numeric version")
		}
		if err := m.Force(v); err != nil {
			return err
		}
		fmt.Printf("forced version to %d\n", v)

	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
	return nil
}

// newMigration creates the next sequentially numbered up/down file pair.
func newMigration(name string) error {
	const dir = "migrations"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading %s (run from backend/): %w", dir, err)
	}

	last := 0
	for _, e := range entries {
		var n int
		if _, err := fmt.Sscanf(e.Name(), "%06d_", &n); err == nil && n > last {
			last = n
		}
	}

	base := fmt.Sprintf("%06d_%s", last+1, name)
	for _, suffix := range []string{".up.sql", ".down.sql"} {
		path := filepath.Join(dir, base+suffix)
		if err := os.WriteFile(path, []byte("-- "+base+suffix+"\n"), 0o644); err != nil {
			return err
		}
		fmt.Println("created", path)
	}
	return nil
}
