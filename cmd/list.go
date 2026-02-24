package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"

	"github.com/odysa/histctl/internal/browser"
	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listJSON   bool
)

var listCmd = &cobra.Command{
	Use:   "list [pattern]",
	Short: "List browser history (non-interactive)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		browsers, err := resolveBrowsers()
		if err != nil {
			return err
		}

		opts := browser.ListOptions{Limit: listLimit}
		if len(args) > 0 {
			re, err := regexp.Compile("(?i)" + args[0])
			if err != nil {
				return fmt.Errorf("invalid regex: %w", err)
			}
			opts.Pattern = re
		}

		ctx := context.Background()
		var all []browser.HistoryEntry
		for _, b := range browsers {
			entries, err := b.List(ctx, opts)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %s: %v\n", b.Name(), err)
				continue
			}
			all = append(all, entries...)
		}

		if listJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(all)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "BROWSER\tURL\tTITLE\tTIME")
		for _, e := range all {
			title := e.Title
			if len(title) > 40 {
				title = title[:39] + "…"
			}
			url := e.URL
			if len(url) > 60 {
				url = url[:59] + "…"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				e.Browser, url, title,
				e.VisitTime.Local().Format("2006-01-02 15:04"))
		}
		return w.Flush()
	},
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "n", 50, "Max entries to display")
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(listCmd)
}
