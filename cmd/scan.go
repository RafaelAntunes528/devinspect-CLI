package cmd

import (
	"encoding/json"
	"path/filepath"

	"devinspect/internal/report"
	"devinspect/internal/scanner"
	"devinspect/internal/stats"

	"github.com/spf13/cobra"
)

var jsonOutput bool

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan a project directory",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) == 1 {
			path = args[0]
		}

		abs, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		s := stats.New(abs)
		rootNode, err := scanner.Scan(abs, s)
		if err != nil {
			return err
		}

		if jsonOutput {
			encoder := json.NewEncoder(cmd.OutOrStdout())
			encoder.SetIndent("", "  ")
			return encoder.Encode(s)
		}

		report.Print(s, rootNode)
		return nil
	},
}

func init() {
	scanCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output JSON report")
	rootCmd.AddCommand(scanCmd)
}
