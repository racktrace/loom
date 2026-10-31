package cmd

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/racktrace/loom/internal/scratch"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "loom",
	Short:        "scratch-first note taking, fast",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if stdinIsPipe() {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			return scratch.AppendBlock(string(data))
		}
		return openEditor()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func stdinIsPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice == 0
}

func openEditor() error {
	p, err := scratch.Path()
	if err != nil {
		return err
	}
	// pastikan file ada supaya editor tidak membuka buffer kosong tanpa nama
	if f, err := os.OpenFile(p, os.O_CREATE, 0o600); err == nil {
		f.Close()
	}

	// $EDITOR bisa berisi argumen, misalnya "code -w", jadi dipecah dulu
	parts := strings.Fields(os.Getenv("EDITOR"))
	if len(parts) == 0 {
		parts = []string{"vi"}
	}
	c := exec.Command(parts[0], append(parts[1:], p)...)
	c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	return c.Run()
}
