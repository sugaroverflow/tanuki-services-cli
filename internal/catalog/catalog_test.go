package catalog

import (
	"testing"
)

func TestFindByName(t *testing.T) {
	svcs := []Service{
		{Name: "payments-api", Team: "payments"},
		{Name: "auth-service", Team: "platform"},
	}
	if s := FindByName(svcs, "auth-service"); s == nil || s.Team != "platform" {
		t.Fatalf("FindByName(auth-service) = %v, want platform", s)
	}
	if s := FindByName(svcs, "missing"); s != nil {
		t.Fatalf("FindByName(missing) = %v, want nil", s)
	}
}

func TestFilterByTeam(t *testing.T) {
	svcs := []Service{
		{Name: "a", Team: "platform"},
		{Name: "b", Team: "payments"},
		{Name: "c", Team: "platform"},
	}
	out := FilterByTeam(svcs, "platform")
	if len(out) != 2 {
		t.Fatalf("FilterByTeam(platform) len = %d, want 2", len(out))
	}
	if out[0].Name != "a" || out[1].Name != "c" {
		t.Fatalf("FilterByTeam(platform) = %v", out)
	}
}
