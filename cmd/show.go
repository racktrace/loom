package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/racktrace/loom/internal/config"
	"github.com/racktrace/loom/internal/store"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:               "show [name|scratch]",
	Short:             "List notes, or print a note or the scratch",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completeNoteNames(true),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		s := store.New(cfg.DataDir)

		// loom show
		if len(args) == 0 {
			return listNotes(s)
		}

		name := args[0]

		// loom show scratch
		if name == store.ScratchName {
			err := printFile(s.ScratchPath())
			if errors.Is(err, fs.ErrNotExist) {
				return errors.New("scratch is empty, run `loom` to start writing")
			}
			return err
		}

		// loom show <name>
		p, err := s.NotePath(name)
		if err != nil {
			return err
		}
		err = printFile(p)
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("note %q does not exist, use `loom add %s`", name, name)
		}
		return err
	},
}

func listNotes(s *store.Store) error {
	names, err := s.ListNotes()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Println("no notes yet, create one with `loom add <name>`")
		return nil
	}
	for _, n := range names {
		fmt.Println(n)
	}
	return nil
}

// printFile streams a file to stdout without loading it fully into memory.
func printFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(os.Stdout, f)
	return err
}

func init() {
	rootCmd.AddCommand(showCmd)
}
