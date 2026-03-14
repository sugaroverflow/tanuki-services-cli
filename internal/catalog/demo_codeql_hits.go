package catalog

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// Demo-only handlers to create reliable SAST/CodeQL findings.

// DemoSQLSearch shows SQL built from request input.
func DemoSQLSearch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	team := r.URL.Query().Get("team")
	q := fmt.Sprintf("SELECT name FROM services WHERE team = '%s'", team)
	_, _ = db.Query(q)
	_, _ = w.Write([]byte("ok"))
}

// DemoFileRead shows path construction from request input.
func DemoFileRead(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("file")
	p := filepath.Join("registry", name)
	_, _ = os.ReadFile(p)
	_, _ = w.Write([]byte("ok"))
}

// DemoShellRun shows shell execution from request input.
func DemoShellRun(w http.ResponseWriter, r *http.Request) {
	arg := r.URL.Query().Get("q")
	_ = exec.Command("sh", "-c", "echo "+arg).Run()
	_, _ = w.Write([]byte("ok"))
}
