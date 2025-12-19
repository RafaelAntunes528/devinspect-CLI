package scanner

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"devinspect/internal/stats"
)

var extensionToLanguage = map[string]string{
	".go":  "Go",
	".js":  "JavaScript",
	".ts":  "TypeScript",
	".py":  "Python",
	".tsx": "TypeScriptReact",
	".jsx": "JavaScriptReact",
}

// Scan now returns the root of the file tree along with stats
func Scan(root string, s *stats.Stats) (*stats.FileNode, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	ignore := loadIgnoreFile(root)

	// WalkDir for counting files and gathering stats
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && shouldSkipDir(d.Name(), ignore) {
			return filepath.SkipDir
		}

		if !d.IsDir() {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()

				lines := countLines(p)
				ext := filepath.Ext(p)

				mu.Lock()
				s.FileCount++
				s.LineCount += lines
				if lang, ok := extensionToLanguage[ext]; ok {
					s.Languages[lang]++
				}
				checkSpecialFiles(filepath.Base(p), s)
				mu.Unlock()
			}(path)
		}

		return nil
	})

	wg.Wait() // Wait for all goroutines

	// Build tree for printing
	rootNode, _ := stats.BuildTree(root, ignore)

	return rootNode, err
}

func shouldSkipDir(name string, ignore []string) bool {
	if name == ".git" || name == "node_modules" {
		return true
	}
	for _, v := range ignore {
		if name == v {
			return true
		}
	}
	return false
}

func loadIgnoreFile(root string) []string {
	path := filepath.Join(root, ".devinspectignore")
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var ignores []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ignores = append(ignores, strings.TrimSpace(scanner.Text()))
	}
	return ignores
}

func checkSpecialFiles(name string, s *stats.Stats) {
	lower := strings.ToLower(name)
	if lower == "readme.md" {
		s.HasReadme = true
	}
	if lower == ".gitignore" {
		s.HasGitignore = true
	}
}

func countLines(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
