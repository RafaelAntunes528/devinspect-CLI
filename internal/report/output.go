package report

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"devinspect/internal/stats"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func Print(s *stats.Stats, rootNode *stats.FileNode) {
	color.Cyan("\n📂 Project Tree:")
	stats.PrintTree(rootNode, "", true)
	fmt.Println()

	calculateScore(s)

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Metric", "Value"})

	table.Append([]string{"Files", fmt.Sprint(s.FileCount)})
	table.Append([]string{"Lines", fmt.Sprint(s.LineCount)})
	table.Append([]string{"Score", fmt.Sprintf("%d / 100", s.Score)})

	if s.HasReadme {
		table.Append([]string{"README.md", color.GreenString("✅ Present")})
	} else {
		table.Append([]string{"README.md", color.RedString("❌ Missing")})
	}

	if s.HasGitignore {
		table.Append([]string{".gitignore", color.GreenString("✅ Present")})
	} else {
		table.Append([]string{".gitignore", color.RedString("❌ Missing")})
	}

	table.Render()

	// Languages
	color.Cyan("\n💻 Languages:")
	for lang, count := range s.Languages {
		bar := strings.Repeat("▓", count)
		fmt.Printf("%-15s %s %d files\n", lang, bar, count)
	}

	// Largest files
	color.Magenta("\n📌 Largest Files:")
	for _, f := range stats.TopNLargestFiles(rootNode, 5) {
		fmt.Printf(" - %s (%d lines, %s, C:%d)\n", f.Path, f.Lines, humanSize(f.Size), f.Complexity)
	}

	// Dependencies
	if len(s.Dependencies) > 0 {
		PrintSortedDependencies(s.Dependencies)
	}

	PrintSuggestions(s)
}

func calculateScore(s *stats.Stats) {
	score := 30
	if s.HasReadme {
		score += 20
	}
	if s.HasGitignore {
		score += 20
	}
	if len(s.Languages) > 1 {
		score += 10
	}
	if s.FileCount > 10 {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	s.Score = score
}

func PrintSuggestions(s *stats.Stats) {
	color.Yellow("\n💡 Suggestions to improve your score:")
	if !s.HasReadme {
		fmt.Println(" - Add a README.md file to explain your project.")
	}
	if !s.HasGitignore {
		fmt.Println(" - Add a .gitignore file to ignore unnecessary files.")
	}
	if len(s.Languages) <= 1 {
		fmt.Println(" - Use multiple languages or frameworks appropriately.")
	}
	if s.FileCount <= 10 {
		fmt.Println(" - Add more meaningful files/modules to expand your project.")
	}
	if s.Score == 100 {
		color.Green("✅ Your project is perfect! No improvements needed.")
	}
}

func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func PrintSortedDependencies(deps map[string]int) {
	color.Cyan("\n📦 Dependencies:")

	// Convert map to slice
	type depCount struct {
		Name  string
		Count int
	}
	var depList []depCount
	for name, count := range deps {
		depList = append(depList, depCount{name, count})
	}

	// Sort descending by count
	sort.Slice(depList, func(i, j int) bool {
		return depList[i].Count > depList[j].Count
	})

	// Print
	for _, d := range depList {
		fmt.Printf(" - %s (%d imports)\n", d.Name, d.Count)
	}
}
