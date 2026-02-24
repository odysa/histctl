package cmd

import (
	"fmt"
	"os"

	"github.com/odysa/histctl/internal/browser"
	"github.com/odysa/histctl/internal/tui"
	"github.com/spf13/cobra"
)

var browserFlag string

var rootCmd = &cobra.Command{
	Use:   "histctl",
	Short: "Browser history manager with regex search and beautiful TUI",
	Long:  "Search, visualize, and delete browser history across Safari, Chrome, Edge, and Firefox.",
	RunE: func(cmd *cobra.Command, args []string) error {
		browsers, err := resolveBrowsers()
		if err != nil {
			return err
		}
		if len(browsers) == 0 {
			fmt.Fprintln(os.Stderr, "No supported browsers found.")
			return nil
		}
		return tui.Run(browsers)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&browserFlag, "browser", "b", "all",
		"Target browser: safari|chrome|edge|firefox|all")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func resolveBrowsers() ([]browser.Browser, error) {
	if browserFlag == "all" {
		return browser.Available(), nil
	}
	b, err := browser.Get(browserFlag)
	if err != nil {
		return nil, err
	}
	if _, err := b.DBPath(); err != nil {
		return nil, fmt.Errorf("%s is not installed or history not found: %w", browserFlag, err)
	}
	return []browser.Browser{b}, nil
}
