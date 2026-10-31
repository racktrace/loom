package cmd

import (
	"strings"

	"github.com/racktrace/loom/internal/scratch"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add <teks>",
	Aliases: []string{"a"},
	Short:   "Tambah satu baris ke scratch",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return scratch.Append(strings.Join(args, " "))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
