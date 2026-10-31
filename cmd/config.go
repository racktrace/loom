package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/racktrace/loom/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the active configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		p, err := config.Path()
		if err != nil {
			return err
		}
		status := "not found, using defaults"
		if _, err := os.Stat(p); err == nil {
			status = "loaded"
		}

		fmt.Printf("config file   : %s (%s)\n", p, status)
		fmt.Printf("data_dir      : %s\n", cfg.DataDir)
		fmt.Printf("editor        : %s\n", strings.Join(cfg.EditorCommand(), " "))
		fmt.Printf("date_format   : %s\n", cfg.Scratch.DateFormat)
		fmt.Printf("default_tags  : %v\n", cfg.Note.DefaultTags)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
