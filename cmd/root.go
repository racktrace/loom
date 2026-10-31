package cmd

import (
	"os"
	"time"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/editor"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "loom",
	Short:        "Scratch-first note taking for the terminal",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := store.New(cfg.DataDir)

		header, err := s.BeginScratchSession(cfg.Scratch.DateFormat, time.Now())
		if err != nil {
			return err
		}

		editErr := editor.Open(cfg.EditorCommand(), s.ScratchPath())

		// Always clean up an empty header, even if the editor failed.
		if err := s.EndScratchSession(header); err != nil {
			return err
		}
		return editErr
	},
}

// Execute runs the root command and exits with status 1 on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
