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

func Scan(root string, s *stats.Stats) (*stats.FileNode, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	depsMap := make(map[string]int)

	ignore := loadIgnoreFile(root)

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

				deps := countDependencies(p)

				mu.Lock()
				s.FileCount++
				s.LineCount += lines
				if lang, ok := extensionToLanguage[ext]; ok {
					s.Languages[lang]++
				}
				checkSpecialFiles(filepath.Base(p), s)

				// Merge dependencies
				for d, c := range deps {
					depsMap[d] += c
				}
				mu.Unlock()
			}(path)
		}
		return nil
	})

	wg.Wait()
	rootNode, _ := stats.BuildTree(root, ignore)
	attachComplexity(rootNode)
	s.Dependencies = depsMap
	return rootNode, err
}

func calculateComplexity(path string) int {
	file, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	complexity := 1
	keywords := []string{"if", "for", "while", "case", "&&", "||", "switch"}
	for scanner.Scan() {
		line := scanner.Text()
		for _, k := range keywords {
			if strings.Contains(line, k) {
				complexity++
			}
		}
	}
	return complexity
}

func countDependencies(path string) map[string]int {
	deps := make(map[string]int)
	file, err := os.Open(path)
	if err != nil {
		return deps
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inImportBlock := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Detect Go multi-line import block
		if strings.HasPrefix(line, "import (") {
			inImportBlock = true
			continue
		}
		if inImportBlock {
			if line == ")" {
				inImportBlock = false
				continue
			}
			dep := strings.Trim(line, "\"")
			if dep != "" {
				deps[dep]++
			}
			continue
		}

		// Single-line import (Go/JS)
		if strings.HasPrefix(line, "import ") || strings.HasPrefix(line, "require(") {
			// Extract only quoted part
			start := strings.Index(line, "\"")
			end := strings.LastIndex(line, "\"")
			if start >= 0 && end > start {
				dep := line[start+1 : end]
				deps[dep]++
			}
		}
	}
	return deps
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

func checkSpecialFiles(name string, s *stats.Stats) {
	lower := strings.ToLower(name)
	if lower == "readme.md" {
		s.HasReadme = true
	}
	if lower == ".gitignore" {
		s.HasGitignore = true
	}
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

func attachComplexity(node *stats.FileNode) {
	if !node.IsDir {
		node.Complexity = calculateComplexity(node.Path)
	}
	for _, child := range node.Children {
		attachComplexity(child)
	}
}
