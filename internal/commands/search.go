package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ar4mirez/aicof/internal/core"
	"github.com/ar4mirez/aicof/internal/ui"
	"github.com/spf13/cobra"
)

// SearchResult represents a search match
type SearchResult struct {
	Name        string
	Type        string
	Description string
	Score       int
	Installed   bool
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for components by keyword",
	Long: `Search for AICoF components across languages, frameworks, and workflows.

Supports fuzzy matching for typos and partial matches. Results are sorted by relevance.

Examples:
  aicof search react              # Search all component types
  aicof search --type fw django   # Search only frameworks
  aicof search py                 # Fuzzy match finds "python"
  aicof search "spring boot"      # Multi-word search

Types (with aliases):
  language   (lang, l)   Language guides
  framework  (fw, f)     Framework guides
  workflow   (wf, w)     Workflow templates`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringP("type", "t", "", "Filter by type: language/lang/l, framework/fw/f, workflow/wf/w")
	searchCmd.Flags().IntP("limit", "l", 20, "Limit number of results")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.ToLower(args[0])
	typeFilter, _ := cmd.Flags().GetString("type")
	limit, _ := cmd.Flags().GetInt("limit")

	typeFilter = normalizeTypeFilter(typeFilter)
	config, _ := core.LoadConfig()
	results := searchComponents(query, typeFilter, config)

	if len(results) == 0 {
		ui.Warn("No components found matching '%s'", query)
		ui.Info("Try a different search term or use 'aicof list --available' to see all components")
		return nil
	}

	results = sortAndLimitResults(results, limit)
	displaySearchResults(query, results)
	return nil
}

func sortAndLimitResults(results []SearchResult, limit int) []SearchResult {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > limit {
		return results[:limit]
	}
	return results
}

func displaySearchResults(query string, results []SearchResult) {
	ui.Bold("Search results for '%s'", query)
	fmt.Println()

	languages := filterByType(results, "language")
	frameworks := filterByType(results, "framework")
	workflows := filterByType(results, "workflow")

	if len(languages) > 0 {
		ui.Section(fmt.Sprintf("Languages (%d)", len(languages)))
		displayResults(languages)
	}
	if len(frameworks) > 0 {
		ui.Section(fmt.Sprintf("Frameworks (%d)", len(frameworks)))
		displayResults(frameworks)
	}
	if len(workflows) > 0 {
		ui.Section(fmt.Sprintf("Workflows (%d)", len(workflows)))
		displayResults(workflows)
	}

	fmt.Println()
	ui.Dim("%d result(s) found", len(results))
}

func searchComponents(query, typeFilter string, config *core.Config) []SearchResult {
	var results []SearchResult

	if typeFilter == "" || typeFilter == "language" {
		results = append(results, searchLanguages(query, config)...)
	}
	if typeFilter == "" || typeFilter == "framework" {
		results = append(results, searchFrameworks(query, config)...)
	}
	if typeFilter == "" || typeFilter == "workflow" {
		results = append(results, searchWorkflows(query, config)...)
	}

	return results
}

func searchLanguages(query string, config *core.Config) []SearchResult {
	var results []SearchResult
	for _, lang := range core.Languages {
		if score := matchScore(query, lang.Name, lang.Description); score > 0 {
			results = append(results, SearchResult{
				Name:        lang.Name,
				Type:        "language",
				Description: lang.Description,
				Score:       score,
				Installed:   config != nil && config.HasLanguage(lang.Name),
			})
		}
	}
	return results
}

func searchFrameworks(query string, config *core.Config) []SearchResult {
	var results []SearchResult
	for _, fw := range core.Frameworks {
		if score := matchScore(query, fw.Name, fw.Description); score > 0 {
			results = append(results, SearchResult{
				Name:        fw.Name,
				Type:        "framework",
				Description: fw.Description,
				Score:       score,
				Installed:   config != nil && config.HasFramework(fw.Name),
			})
		}
	}
	return results
}

func searchWorkflows(query string, config *core.Config) []SearchResult {
	var results []SearchResult
	for _, wf := range core.Workflows {
		if score := matchScore(query, wf.Name, wf.Description); score > 0 {
			results = append(results, SearchResult{
				Name:        wf.Name,
				Type:        "workflow",
				Description: wf.Description,
				Score:       score,
				Installed:   config != nil && config.HasWorkflow(wf.Name),
			})
		}
	}
	return results
}

// matchScore returns a score indicating match quality (0 = no match)
func matchScore(query, name, description string) int {
	queryLower := strings.ToLower(query)
	nameLower := strings.ToLower(name)
	descLower := strings.ToLower(description)

	score := 0

	// Exact match on name (highest score)
	if nameLower == queryLower {
		return 100
	}

	// Prefix match on name
	if strings.HasPrefix(nameLower, queryLower) {
		score = 80
	}

	// Contains match on name
	if score == 0 && strings.Contains(nameLower, queryLower) {
		score = 60
	}

	// Contains match on description
	if score == 0 && strings.Contains(descLower, queryLower) {
		score = 40
	}

	// Fuzzy match on name (Levenshtein distance)
	if score == 0 {
		dist := levenshteinDistance(queryLower, nameLower)
		if dist <= 2 && dist < len(nameLower)/2 {
			score = 30 - dist*5
		}
	}

	return score
}

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func normalizeTypeFilter(filter string) string {
	switch strings.ToLower(filter) {
	case "language", "lang", "l":
		return "language"
	case "framework", "fw", "f":
		return "framework"
	case "workflow", "wf", "w":
		return "workflow"
	default:
		return ""
	}
}

func filterByType(results []SearchResult, componentType string) []SearchResult {
	var filtered []SearchResult
	for _, r := range results {
		if r.Type == componentType {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func displayResults(results []SearchResult) {
	for _, r := range results {
		if r.Installed {
			ui.SuccessItem(1, "%s - %s (installed)", r.Name, r.Description)
		} else {
			ui.ListItem(1, "%s - %s", r.Name, r.Description)
		}
	}
}
