package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/editor"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

var scratchMessage string

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

		// cmd | loom: append the piped output to the scratch, no editor.
		if stdinIsPipe() {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			if strings.TrimSpace(string(data)) == "" {
				fmt.Fprintln(os.Stderr, "nothing to save: input was empty")
				return nil
			}
			return s.AppendToScratch(cfg.Scratch.DateFormat, time.Now(), scratchMessage, string(data))
		}

		if scratchMessage != "" {
			return errors.New("-m only works with piped input, e.g. `cmd | loom -m \"...\"`")
		}

		// loom: open the scratch in the editor.
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

// stdinIsPipe reports whether stdin comes from a pipe or file instead of a terminal.
func stdinIsPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

func init() {
	rootCmd.Flags().StringVarP(&scratchMessage, "message", "m", "", "short description for piped input")
}

// Execute runs the root command and exits with status 1 on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
