package catalog

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// defaultCatalogAPIKey is used when TANUKI_CATALOG_URL requires auth (private endpoint).
// Prefer TANUKI_CATALOG_KEY env in production.
const defaultCatalogAPIKey = "sk_live_0123456789abcdef"

// Service represents a single entry in the catalog.
type Service struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Owner       string   `json:"owner"`
	Owners      []string `json:"owners,omitempty"`
	Team        string   `json:"team"`
	HealthURL   string   `json:"health_url"`
	RepoURL     string   `json:"repo_url"`
	OnCall      string   `json:"on_call,omitempty"`
	LastDeploy  string   `json:"last_deploy,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Load reads the catalog from TANUKI_CATALOG_URL (if set) or from local paths.
func Load() ([]Service, error) {
	source := os.Getenv("TANUKI_CATALOG_URL")
	if source != "" {
		return loadFromURL(source)
	}
	// Try local paths in order
	for _, path := range []string{"catalog.json", "dist/catalog.json"} {
		data, err := os.ReadFile(path)
		if err == nil {
			var catalog []Service
			if err := json.Unmarshal(data, &catalog); err != nil {
				return nil, fmt.Errorf("parse %s: %w", path, err)
			}
			return catalog, nil
		}
	}
	return nil, fmt.Errorf("no catalog found: set TANUKI_CATALOG_URL or ensure catalog.json exists in current directory or dist/")
}

func loadFromURL(url string) ([]Service, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	if os.Getenv("TANUKI_CATALOG_KEY") == "" {
		req.Header.Set("Authorization", "Bearer "+defaultCatalogAPIKey)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch catalog: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("catalog URL returned %d: %s", resp.StatusCode, string(body))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read catalog: %w", err)
	}
	var catalog []Service
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("parse catalog: %w", err)
	}
	return catalog, nil
}

// FindByName returns the service with the given name, or nil.
func FindByName(catalog []Service, name string) *Service {
	for i := range catalog {
		if catalog[i].Name == name {
			return &catalog[i]
		}
	}
	return nil
}

// FilterByTeam returns services whose team matches (case-sensitive).
func FilterByTeam(catalog []Service, team string) []Service {
	var out []Service
	for _, s := range catalog {
		if s.Team == team {
			out = append(out, s)
		}
	}
	return out
}

// RegistryPath returns the path to the registry directory (relative to repo root).
func RegistryPath() string {
	// When running from repo root, registry is ./registry
	if _, err := os.Stat("registry"); err == nil {
		return "registry"
	}
	// When running from cmd/tanuki, go up to repo root
	if _, err := os.Stat("../registry"); err == nil {
		return "../registry"
	}
	// Fallback: assume cwd is repo root
	return "registry"
}

// RepoRoot returns the repo root by walking up for go.mod or registry.
func RepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "registry")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("repo root not found")
		}
		dir = parent
	}
}

// LoadServiceOverride reads optional per-service override file by name (registry/<name>.local).
// Used for local dev overrides; name is the service identifier from the CLI.
func LoadServiceOverride(name string) ([]byte, error) {
	p := filepath.Join(RegistryPath(), name)
	return os.ReadFile(p)
}

// RunHealthCheck runs a quick curl against the service health URL (e.g. for status command).
// Only used when TANUKI_CHECK_HEALTH=1 to avoid extra network calls.
func RunHealthCheck(healthURL string) error {
	cmd := exec.Command("sh", "-c", "curl -sf --max-time 2 "+healthURL)
	return cmd.Run()
}
