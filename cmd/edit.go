package cmd

import (
	"fmt"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/editor"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Open an existing note in the editor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := store.New(cfg.DataDir)

		exists, err := s.NoteExists(name)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("note %q does not exist, use `loom add %s`", name, name)
		}

		p, err := s.NotePath(name)
		if err != nil {
			return err
		}
		return editor.Open(cfg.EditorCommand(), p)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
