package commands

import (
	"testing"

	"github.com/ar4mirez/aicof/internal/core"
)

func TestMatchScore(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		compName    string
		description string
		wantScore   int
	}{
		{
			name:        "exact match",
			query:       "react",
			compName:    "react",
			description: "React framework",
			wantScore:   100,
		},
		{
			name:        "exact match case insensitive",
			query:       "React",
			compName:    "react",
			description: "React framework",
			wantScore:   100,
		},
		{
			name:        "prefix match",
			query:       "type",
			compName:    "typescript",
			description: "TypeScript language",
			wantScore:   80,
		},
		{
			name:        "contains match in name",
			query:       "script",
			compName:    "typescript",
			description: "TypeScript language",
			wantScore:   60,
		},
		{
			name:        "contains match in description",
			query:       "framework",
			compName:    "react",
			description: "React framework",
			wantScore:   40,
		},
		{
			name:        "fuzzy match - one char off",
			query:       "pythn",
			compName:    "python",
			description: "Python language",
			wantScore:   25, // 30 - 1*5
		},
		{
			name:        "prefix match - pytho",
			query:       "pytho",
			compName:    "python",
			description: "Python language",
			wantScore:   80, // prefix match takes precedence
		},
		{
			name:        "no match",
			query:       "xyz123",
			compName:    "react",
			description: "React framework",
			wantScore:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchScore(tt.query, tt.compName, tt.description)
			if got != tt.wantScore {
				t.Errorf("matchScore(%q, %q, %q) = %d, want %d",
					tt.query, tt.compName, tt.description, got, tt.wantScore)
			}
		})
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		name string
		s1   string
		s2   string
		want int
	}{
		{"identical", "hello", "hello", 0},
		{"one insertion", "hello", "helo", 1},
		{"one deletion", "helo", "hello", 1},
		{"one substitution", "hello", "hallo", 1},
		{"empty first", "", "hello", 5},
		{"empty second", "hello", "", 5},
		{"both empty", "", "", 0},
		{"completely different", "abc", "xyz", 3},
		{"case sensitive", "Hello", "hello", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := levenshteinDistance(tt.s1, tt.s2)
			if got != tt.want {
				t.Errorf("levenshteinDistance(%q, %q) = %d, want %d",
					tt.s1, tt.s2, got, tt.want)
			}
		})
	}
}

func TestNormalizeTypeFilter(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"language", "language"},
		{"lang", "language"},
		{"l", "language"},
		{"LANGUAGE", "language"},
		{"framework", "framework"},
		{"fw", "framework"},
		{"f", "framework"},
		{"workflow", "workflow"},
		{"wf", "workflow"},
		{"w", "workflow"},
		{"", ""},
		{"invalid", ""},
		{"xyz", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeTypeFilter(tt.input)
			if got != tt.want {
				t.Errorf("normalizeTypeFilter(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestFilterByType(t *testing.T) {
	results := []SearchResult{
		{Name: "go", Type: "language"},
		{Name: "react", Type: "framework"},
		{Name: "python", Type: "language"},
		{Name: "create-prd", Type: "workflow"},
		{Name: "nextjs", Type: "framework"},
	}

	tests := []struct {
		filterType string
		wantCount  int
	}{
		{"language", 2},
		{"framework", 2},
		{"workflow", 1},
		{"invalid", 0},
	}

	for _, tt := range tests {
		t.Run(tt.filterType, func(t *testing.T) {
			got := filterByType(results, tt.filterType)
			if len(got) != tt.wantCount {
				t.Errorf("filterByType(results, %q) returned %d items, want %d",
					tt.filterType, len(got), tt.wantCount)
			}
		})
	}
}

func TestSortAndLimitResults(t *testing.T) {
	results := []SearchResult{
		{Name: "a", Score: 50},
		{Name: "b", Score: 100},
		{Name: "c", Score: 75},
		{Name: "d", Score: 25},
		{Name: "e", Score: 90},
	}

	tests := []struct {
		name      string
		limit     int
		wantFirst string
		wantCount int
	}{
		{"limit 3", 3, "b", 3},
		{"limit 10", 10, "b", 5},
		{"limit 1", 1, "b", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid mutating original
			resultsCopy := make([]SearchResult, len(results))
			copy(resultsCopy, results)

			got := sortAndLimitResults(resultsCopy, tt.limit)
			if len(got) != tt.wantCount {
				t.Errorf("sortAndLimitResults() returned %d items, want %d", len(got), tt.wantCount)
			}
			if got[0].Name != tt.wantFirst {
				t.Errorf("sortAndLimitResults() first item = %q, want %q", got[0].Name, tt.wantFirst)
			}
		})
	}
}

func TestSearchLanguages(t *testing.T) {
	// Test with nil config
	results := searchLanguages("go", nil)
	found := false
	for _, r := range results {
		if r.Name == "go" {
			found = true
			if r.Installed {
				t.Error("Expected Installed=false with nil config")
			}
		}
	}
	if !found {
		t.Error("Expected to find 'go' in results")
	}
}

func TestSearchFrameworks(t *testing.T) {
	results := searchFrameworks("react", nil)
	found := false
	for _, r := range results {
		if r.Name == "react" {
			found = true
			if r.Type != "framework" {
				t.Errorf("Expected type=framework, got %q", r.Type)
			}
		}
	}
	if !found {
		t.Error("Expected to find 'react' in results")
	}
}

func TestSearchWorkflows(t *testing.T) {
	results := searchWorkflows("prd", nil)
	found := false
	for _, r := range results {
		if r.Name == "create-prd" {
			found = true
			if r.Type != "workflow" {
				t.Errorf("Expected type=workflow, got %q", r.Type)
			}
		}
	}
	if !found {
		t.Error("Expected to find 'create-prd' in results")
	}
}

func TestSearchComponents(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		typeFilter string
		wantMin    int
	}{
		{"search all for go", "go", "", 1},
		{"search languages only", "python", "language", 1},
		{"search frameworks only", "react", "framework", 1},
		{"no results", "xyz123nonexistent", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := searchComponents(tt.query, tt.typeFilter, nil)
			if len(results) < tt.wantMin {
				t.Errorf("searchComponents(%q, %q) returned %d results, want at least %d",
					tt.query, tt.typeFilter, len(results), tt.wantMin)
			}
		})
	}
}

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		Name:        "test",
		Type:        "language",
		Description: "Test language",
		Score:       100,
		Installed:   true,
	}

	if r.Name != "test" {
		t.Errorf("Name = %q, want %q", r.Name, "test")
	}
	if r.Type != "language" {
		t.Errorf("Type = %q, want %q", r.Type, "language")
	}
	if r.Description != "Test language" {
		t.Errorf("Description = %q, want %q", r.Description, "Test language")
	}
	if r.Score != 100 {
		t.Errorf("Score = %d, want %d", r.Score, 100)
	}
	if !r.Installed {
		t.Error("Installed = false, want true")
	}
}

// Test that core.Languages, Frameworks, Workflows are accessible
func TestRegistryAccess(t *testing.T) {
	if len(core.Languages) == 0 {
		t.Error("core.Languages is empty")
	}
	if len(core.Frameworks) == 0 {
		t.Error("core.Frameworks is empty")
	}
	if len(core.Workflows) == 0 {
		t.Error("core.Workflows is empty")
	}
}
