package stats

import "path/filepath"

type Stats struct {
	ProjectName  string         `json:"project_name"`
	RootPath     string         `json:"root_path"`
	FileCount    int            `json:"file_count"`
	LineCount    int            `json:"line_count"`
	Languages    map[string]int `json:"languages"`
	HasReadme    bool           `json:"has_readme"`
	HasGitignore bool           `json:"has_gitignore"`
	Score        int            `json:"score"`
	Dependencies map[string]int `json:"dependencies"`
}

func New(root string) *Stats {
	return &Stats{
		ProjectName: filepath.Base(root),
		RootPath:    root,
		Languages:   make(map[string]int),
	}
}
