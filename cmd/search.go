package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:               `search "<query>"`,
	Short:             "Search the scratch and all notes (case-insensitive)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		if strings.TrimSpace(query) == "" {
			return errors.New("search query cannot be empty")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := store.New(cfg.DataDir)

		matches, err := s.Search(query)
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			fmt.Fprintf(os.Stderr, "no matches for %q\n", query)
			return nil
		}
		for _, m := range matches {
			fmt.Printf("%s:%d: %s\n", m.Source, m.Line, strings.TrimSpace(m.Text))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
