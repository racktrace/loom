package scratch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const timeFmt = "2006-01-02 15:04"

// Dir mengembalikan folder data loom.
// LOOM_DIR bisa dipakai untuk override (berguna untuk testing/benchmark).
func Dir() (string, error) {
	if d := os.Getenv("LOOM_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".loom"), nil
}

// Path mengembalikan lokasi scratch.md, sekaligus memastikan foldernya ada.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "scratch.md"), nil
}

// write membuka file dengan O_APPEND: tidak membaca isi lama sama sekali,
// jadi kecepatannya tetap sama walaupun scratch sudah besar.
func write(s string) error {
	p, err := Path()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(s); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// Append menambah satu baris bertimestamp.
func Append(text string) error {
	return write(fmt.Sprintf("- %s %s\n", time.Now().Format(timeFmt), text))
}

// AppendBlock menambah teks multi-baris (misalnya dari pipe) sebagai code block.
func AppendBlock(text string) error {
	text = strings.TrimRight(text, "\n")
	if text == "" {
		return nil
	}
	return write(fmt.Sprintf("- %s\n```\n%s\n```\n", time.Now().Format(timeFmt), text))
}
