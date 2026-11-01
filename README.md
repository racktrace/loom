# loom

Scratch-first note taking for the terminal.

loom gives you two places to write:

- **scratch**: a single file for dumping thoughts, ideas and brainstorming
  without deciding on a title, a folder or a tag first.
- **notes**: named Markdown files with frontmatter, for things worth keeping.

Both open in your own editor. loom stays out of the way.

## Principles

- **Scratch-first.** Typing `loom` is all it takes to start writing.
- **Fast.** No database, no network, no background process.
- **File-based.** Everything is plain Markdown you can read, grep or sync
  with any tool.

## Install

Build from source (requires Go):

```bash
git clone https://github.com/racktrace/loom.git
cd loom
go build -trimpath -ldflags="-s -w" -o loom .
mkdir -p ~/.local/bin
cp loom ~/.local/bin/
```

Make sure `~/.local/bin` is in your `PATH`.

## Usage

```
loom                  open the scratch in your editor
loom add <name>       create a new note and open it
loom edit <name>      open an existing note
loom show             list all notes
loom show <name>      print a note
loom show scratch     print the scratch
loom config           show the active configuration
loom search "<query>"  search the scratch and all notes (case-insensitive)
cmd | loom            append command output to the scratch
cmd | loom -m "msg"   same, with a short description
```

### Scratch

Every time you run `loom`, a session header is added to the bottom of the
scratch before your editor opens:

```markdown
# Scratch

## 2026-11-01 08:30

an idea about something...

## 2026-11-01 11:20

```

If you close the editor without writing anything, the empty header is removed.
Earlier sessions are never rewritten.

### Notes

`loom add systemd-timer` creates `systemd-timer.md` with frontmatter:

```markdown
---
title: "systemd-timer"
created: 2026-11-01T11:20:00+08:00
tags: []
---
```

Note names may contain letters, digits, `-` and `_`. `scratch` is reserved.
`loom add` never overwrites an existing note, and `loom edit` never creates one.

## Configuration

On first run, loom creates `~/.config/loom/config.toml`
(or `$XDG_CONFIG_HOME/loom/config.toml`):

```toml
# editor = "nano"
# data_dir = "~/.local/share/loom"

[scratch]
date_format = "%Y-%m-%d %H:%M"

[note]
default_tags = []
```

Supported date directives: `%Y %y %m %d %H %M %S %b %a %%`.

Settings are resolved in this order, later wins:

1. built-in defaults
2. `config.toml`
3. environment variables: `LOOM_DIR`, `LOOM_EDITOR`

The editor falls back to `$EDITOR`, then `nano`.

## Data

```
~/.local/share/loom/        # or $XDG_DATA_HOME/loom
├── scratch.md
└── notes/
    └── <name>.md
```

Use `LOOM_DIR` to point loom at another directory, for example while testing:

```bash
LOOM_DIR=/tmp/loom-test loom
```

## Roadmap

- [v] Search across scratch and notes
- [ ] Backup
- [ ] Makefile and release binaries
