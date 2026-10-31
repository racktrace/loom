package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/editor"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Create a new note and open it in the editor",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := store.New(cfg.DataDir)

		p, err := s.CreateNote(name, cfg.Note.DefaultTags, time.Now())
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("note %q already exists, use `loom edit %s`", name, name)
		}
		if err != nil {
			return err
		}

		return editor.Open(cfg.EditorCommand(), p)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
