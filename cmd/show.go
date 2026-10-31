package cmd

import (
	"errors"
	"io"
	"io/fs"
	"os"

	"github.com/racktrace/loom/internal/scratch"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Tampilkan isi scratch",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := scratch.Path()
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if errors.Is(err, fs.ErrNotExist) {
			return nil // scratch belum ada, tidak perlu error
		}
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(os.Stdout, f) // streaming, tidak load seluruh file ke memori
		return err
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
