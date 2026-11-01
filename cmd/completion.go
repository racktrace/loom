package cmd

import (
	"strings"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

// completeNoteNames suggests existing note names for the first argument.
// withScratch also offers "scratch", which only `loom show` accepts.
func completeNoteNames(withScratch bool) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp // only one name is accepted
		}

		cfg, err := config.Load()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names, err := store.New(cfg.DataDir).ListNotes()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		if withScratch {
			names = append([]string{store.ScratchName}, names...)
		}

		var out []string
		for _, n := range names {
			if strings.HasPrefix(n, toComplete) {
				out = append(out, n)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	}
}
