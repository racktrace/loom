package editor

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Open runs the editor command on path and waits until the editor exits.
// command comes from cfg.EditorCommand(), e.g. ["nvim"] or ["code", "-w"].
func Open(command []string, path string) error {
	if len(command) == 0 {
		return errors.New("no editor configured")
	}

	args := append(command[1:], path)
	c := exec.Command(command[0], args...)

	// Hand the terminal over to the editor.
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("editor %q not found: set `editor` in config, LOOM_EDITOR or $EDITOR", command[0])
		}
		return fmt.Errorf("editor exited with error: %w", err)
	}
	return nil
}
