package store

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// ScratchName is reserved so `loom show scratch` is never ambiguous.
const ScratchName = "scratch"

// Store knows where loom keeps its files.
type Store struct {
	dir string
}

// New returns a Store rooted at dir (usually cfg.DataDir).
func New(dir string) *Store {
	return &Store{dir: dir}
}

// Init creates the data directories if they don't exist yet.
func (s *Store) Init() error {
	return os.MkdirAll(s.NotesDir(), 0o700)
}

// ScratchPath returns the path of the single scratch file.
func (s *Store) ScratchPath() string {
	return filepath.Join(s.dir, "scratch.md")
}

// NotesDir returns the directory holding permanent notes.
func (s *Store) NotesDir() string {
	return filepath.Join(s.dir, "notes")
}

// NotePath returns the file path for a note, after validating its name.
func (s *Store) NotePath(name string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}
	return filepath.Join(s.NotesDir(), name+".md"), nil
}

// NoteExists reports whether a note with this name already exists.
func (s *Store) NoteExists(name string) (bool, error) {
	p, err := s.NotePath(name)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(p)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// ValidateName allows letters, digits, '-' and '_' only.
// This also blocks path tricks like "../../etc/passwd".
func ValidateName(name string) error {
	if name == "" {
		return errors.New("note name cannot be empty")
	}
	if strings.EqualFold(name, ScratchName) {
		return fmt.Errorf("%q is a reserved name, choose another one", name)
	}
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return fmt.Errorf("invalid note name %q: use only letters, digits, '-' or '_'", name)
		}
	}
	return nil
}

// BeginScratchSession prepares the scratch file for a new writing session.
// It creates the file with a "# Scratch" title if needed, then appends a
// date header. The header is returned so EndScratchSession can find it.
func (s *Store) BeginScratchSession(dateFormat string, now time.Time) (string, error) {
	if err := s.Init(); err != nil {
		return "", err
	}

	f, err := os.OpenFile(s.ScratchPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", err
	}
	if info.Size() == 0 {
		if _, err := f.WriteString("# Scratch\n"); err != nil {
			return "", err
		}
	}

	header := "## " + FormatDate(dateFormat, now)
	if _, err := f.WriteString("\n" + header + "\n\n"); err != nil {
		return "", err
	}
	return header, f.Close()
}

// EndScratchSession removes the session header if nothing was written under it.
func (s *Store) EndScratchSession(header string) error {
	b, err := os.ReadFile(s.ScratchPath())
	if err != nil {
		return err
	}
	content := strings.TrimRight(string(b), " \t\n")
	if !strings.HasSuffix(content, header) {
		return nil // something was written after the header: keep it
	}
	trimmed := strings.TrimRight(strings.TrimSuffix(content, header), " \t\n")
	return os.WriteFile(s.ScratchPath(), []byte(trimmed+"\n"), 0o600)
}

// FormatDate formats t using strftime-style directives such as %Y-%m-%d %H:%M.
// Supported: %Y %y %m %d %H %M %S %b %a %%. Unknown directives are kept as-is.
func FormatDate(format string, t time.Time) string {
	var b strings.Builder
	for i := 0; i < len(format); i++ {
		if format[i] != '%' || i+1 == len(format) {
			b.WriteByte(format[i])
			continue
		}
		i++
		switch format[i] {
		case 'Y':
			fmt.Fprintf(&b, "%04d", t.Year())
		case 'y':
			fmt.Fprintf(&b, "%02d", t.Year()%100)
		case 'm':
			fmt.Fprintf(&b, "%02d", int(t.Month()))
		case 'd':
			fmt.Fprintf(&b, "%02d", t.Day())
		case 'H':
			fmt.Fprintf(&b, "%02d", t.Hour())
		case 'M':
			fmt.Fprintf(&b, "%02d", t.Minute())
		case 'S':
			fmt.Fprintf(&b, "%02d", t.Second())
		case 'b':
			b.WriteString(t.Format("Jan"))
		case 'a':
			b.WriteString(t.Format("Mon"))
		case '%':
			b.WriteByte('%')
		default:
			b.WriteByte('%')
			b.WriteByte(format[i])
		}
	}
	return b.String()
}

// CreateNote creates a new note with frontmatter and returns its path.
// It fails with an error wrapping fs.ErrExist if the note already exists.
func (s *Store) CreateNote(name string, tags []string, now time.Time) (string, error) {
	p, err := s.NotePath(name)
	if err != nil {
		return "", err
	}
	if err := s.Init(); err != nil {
		return "", err
	}

	// O_EXCL: create only if the file does not exist yet, never overwrite.
	f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	defer f.Close()

	quoted := make([]string, len(tags))
	for i, t := range tags {
		quoted[i] = fmt.Sprintf("%q", t)
	}

	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "title: %q\n", name)
	fmt.Fprintf(&b, "created: %s\n", now.Format(time.RFC3339))
	fmt.Fprintf(&b, "tags: [%s]\n", strings.Join(quoted, ", "))
	b.WriteString("---\n\n")

	if _, err := f.WriteString(b.String()); err != nil {
		return "", err
	}
	return p, f.Close()
}

// ListNotes returns the names of all notes, sorted alphabetically.
// A missing notes directory simply means there are no notes yet.
func (s *Store) ListNotes() ([]string, error) {
	entries, err := os.ReadDir(s.NotesDir())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var names []string
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".md" {
			continue
		}
		names = append(names, strings.TrimSuffix(e.Name(), ".md"))
	}
	return names, nil
}
