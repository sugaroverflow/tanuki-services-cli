package formatter

import (
	"fmt"
	"strings"

	"github.com/example/tanuki/internal/catalog"
)

// List prints a table of services (name, version, team, owner).
func List(services []catalog.Service) {
	if len(services) == 0 {
		fmt.Println("No services found.")
		return
	}
	// Simple column layout
	fmt.Printf("%-20s %-12s %-12s %s\n", "NAME", "VERSION", "TEAM", "OWNER")
	fmt.Println(strings.Repeat("-", 60))
	for _, s := range services {
		fmt.Printf("%-20s %-12s %-12s %s\n", s.Name, s.Version, s.Team, s.Owner)
	}
}

// Status prints detailed status for one service.
func Status(s *catalog.Service) {
	if s == nil {
		fmt.Println("Service not found.")
		return
	}
	fmt.Printf("Name:        %s\n", s.Name)
	fmt.Printf("Version:     %s\n", s.Version)
	fmt.Printf("Owner:       %s\n", s.Owner)
	fmt.Printf("Team:        %s\n", s.Team)
	fmt.Printf("Health URL:  %s\n", s.HealthURL)
	fmt.Printf("Repo:        %s\n", s.RepoURL)
	fmt.Printf("Last deploy: %s\n", s.LastDeploy)
	if s.Description != "" {
		fmt.Printf("Description: %s\n", s.Description)
	}
}

// Owners prints owner and on-call info for one service.
func Owners(s *catalog.Service) {
	if s == nil {
		fmt.Println("Service not found.")
		return
	}
	fmt.Printf("Service:  %s\n", s.Name)
	fmt.Printf("Owner:    %s\n", s.Owner)
	if len(s.Owners) > 0 {
		fmt.Printf("Others:   %s\n", strings.Join(s.Owners, ", "))
	}
	fmt.Printf("On-call:  %s\n", s.OnCall)
	if s.OnCall == "" {
		fmt.Println("(no on-call info)")
	}
}
