package catalog

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Demo-only patterns for CodeQL/SAST to detect. Do not use in production.

// Hardcoded credential: variable name and literal match go/hardcoded-credentials.
var password = "admin123"

// GetCatalogPassword returns the catalog API password (demo: hardcoded).
func GetCatalogPassword() string {
	return password
}

// ReadUserFile reads a file by name under baseDir. Path injection if name is user-controlled.
func ReadUserFile(baseDir, name string) ([]byte, error) {
	p := filepath.Join(baseDir, name)
	return os.ReadFile(p)
}

// RunShellCommand runs a shell command built from user input (command injection).
func RunShellCommand(userURL string) error {
	cmd := exec.Command("sh", "-c", "curl -sf "+userURL)
	return cmd.Run()
}
