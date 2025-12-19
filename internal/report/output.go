package report

import (
	"fmt"
	"os"
	"strings"

	"devinspect/internal/stats"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func Print(s *stats.Stats, rootNode *stats.FileNode) {
	color.Cyan("\n📂 Project Tree:")
	stats.PrintTree(rootNode, "", true)

	calculateScore(s)

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Metric", "Value"})

	// Metrics
	table.Append([]string{"Files", fmt.Sprint(s.FileCount)})
	table.Append([]string{"Lines", fmt.Sprint(s.LineCount)})
	table.Append([]string{"Score", fmt.Sprintf("%d / 100", s.Score)})

	// Special files
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

	PrintSuggestions(s)
}

func calculateScore(s *stats.Stats) {
	score := 40
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
		fmt.Println(" - Use multiple languages or frameworks appropriately to increase versatility.")
	}
	if s.FileCount <= 10 {
		fmt.Println(" - Add more meaningful files/modules to expand your project.")
	}

	if s.Score == 100 {
		color.Green("✅ Your project is perfect! No improvements needed.")
	}
}
