package stats

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type FileNode struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	IsDir    bool        `json:"is_dir"`
	Lines    int         `json:"lines,omitempty"`
	Size     int64       `json:"size,omitempty"`
	Children []*FileNode `json:"children,omitempty"`
}

func BuildTree(root string, ignore []string) (*FileNode, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	node := &FileNode{
		Name:  info.Name(),
		Path:  root,
		IsDir: info.IsDir(),
		Size:  info.Size(),
	}

	if info.IsDir() {
		entries, _ := os.ReadDir(root)
		for _, e := range entries {
			if shouldSkipDir(e.Name(), ignore) {
				continue
			}
			child, _ := BuildTree(filepath.Join(root, e.Name()), ignore)
			node.Children = append(node.Children, child)
		}
	} else {
		node.Lines = countLines(root)
	}

	return node, nil
}

func PrintTree(node *FileNode, prefix string, last bool) {
	folderColor := color.New(color.FgCyan).SprintFunc()
	fileColor := color.New(color.FgWhite).SprintFunc()
	readmeColor := color.New(color.FgGreen).SprintFunc()
	gitignoreColor := color.New(color.FgYellow).SprintFunc()

	connector := "├── "
	nextPrefix := prefix + "│   "
	if last {
		connector = "└── "
		nextPrefix = prefix + "    "
	}

	name := node.Name
	if node.IsDir {
		fmt.Println(prefix + connector + "📁 " + folderColor(name) + "/")
	} else {
		displayName := "📄 " + fileColor(name)
		if name == "README.md" {
			displayName = "✅ " + readmeColor(name)
		} else if name == ".gitignore" {
			displayName = "⚠️ " + gitignoreColor(name)
		}
		fmt.Printf("%s%s %-20s %4d lines | %6s\n", prefix, connector, displayName, node.Lines, humanSize(node.Size))
	}

	for i, child := range node.Children {
		PrintTree(child, nextPrefix, i == len(node.Children)-1)
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
