package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/odysa/histctl/internal/backup"
	"github.com/odysa/histctl/internal/browser"
	"github.com/odysa/histctl/internal/process"
	"github.com/spf13/cobra"
)

var (
	deleteDryRun  bool
	deleteYes     bool
	deleteNoBackup bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <pattern>",
	Short: "Delete history entries matching a regex pattern",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern, err := regexp.Compile("(?i)" + args[0])
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		browsers, err := resolveBrowsers()
		if err != nil {
			return err
		}

		ctx := context.Background()
		var hadErrors bool

		for _, b := range browsers {
			// Check if running
			running, err := process.IsRunning(b.ProcessName())
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not check if %s is running: %v\n", b.Name(), err)
				continue
			}
			if running {
				fmt.Fprintf(os.Stderr, "error: %s is running — close it first\n", b.Name())
				hadErrors = true
				continue
			}

			if deleteDryRun {
				result, err := b.Delete(ctx, pattern, true)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s: %v\n", b.Name(), err)
					hadErrors = true
					continue
				}
				fmt.Printf("[%s] would delete %d entries\n", b.Name(), result.Matched)

				// Show matching entries
				entries, _ := b.List(ctx, browser.ListOptions{Pattern: pattern, Limit: 20})
				for _, e := range entries {
					fmt.Printf("  %s  %s\n", e.URL, e.VisitTime.Local().Format("2006-01-02 15:04"))
				}
				if result.Matched > 20 {
					fmt.Printf("  ... and %d more\n", result.Matched-20)
				}
				continue
			}

			// Preview count
			result, err := b.Delete(ctx, pattern, true)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s: %v\n", b.Name(), err)
				hadErrors = true
				continue
			}
			if result.Matched == 0 {
				fmt.Printf("[%s] no matching entries\n", b.Name())
				continue
			}

			// Confirm
			if !deleteYes {
				fmt.Printf("[%s] delete %d entries? (y/N): ", b.Name(), result.Matched)
				var answer string
				fmt.Scanln(&answer)
				if answer != "y" && answer != "Y" {
					fmt.Println("  skipped")
					continue
				}
			}

			// Backup
			if !deleteNoBackup {
				dbPath, _ := b.DBPath()
				backupPath, err := backup.Create(dbPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: backup failed for %s: %v\n", b.Name(), err)
					hadErrors = true
					continue
				}
				fmt.Printf("[%s] backed up to %s\n", b.Name(), backupPath)
			}

			// Delete
			result, err = b.Delete(ctx, pattern, false)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s: %v\n", b.Name(), err)
				hadErrors = true
				continue
			}
			fmt.Printf("[%s] deleted %d entries\n", b.Name(), result.Deleted)
		}

		if hadErrors {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteDryRun, "dry-run", "d", false, "Preview matches without deleting")
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "Skip confirmation prompt")
	deleteCmd.Flags().BoolVar(&deleteNoBackup, "no-backup", false, "Skip creating a backup")
	rootCmd.AddCommand(deleteCmd)
}
