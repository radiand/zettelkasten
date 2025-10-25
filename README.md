_zettelkasten_ – creating and managing plain text notes. 100% LLM-free code.

# Overview

Single note is just a header and text:

````
```toml
title = "This is your note"
timestamp = "2024-09-29T11:18:38+02:00"
uid = "20240929T091838Z"
tags = ["lang:en", "topic:demo", "whatever"]
referred_from = []
refers_to = []
```

...and here you write whatever you want. Markdown is encouraged, as all notes
are kept as .md files.

Notes are linked using their UIDs, like [this](20240101T010203Z) or
[[20240202T010203Z]]. Zettelkasten will automatically detect links and update
headers.

If your editor supports it, put ![image](./from/resources/img.jpg).
````

`zettelkasten` command line application comes with commands:

- `$ zettelkasten init` to set things up for the first time,
- `$ zettelkasten new` to create new command,
- `$ zettelkasten link` to find references between notes and fill
  `referred_from`, `refers_to` fields of the header,
- `$ zettelkasten commit` to `git commit` if you keep your notes
  version-controlled.

# Try yourself

## Install

Currently, to install `zettelkasten` you need golang compiler. You distro likely
provides it in its repos.

After cloning this repository, `cd` into and run:

```bash
go install ./cmd/zettelkasten
```

If your OS is setup properly and your `PATH` is populated with `go` built
binaries, you can use `zettelkasten` in your terminal now.

## Initialize

To follow default settings and create configuration in
`$HOME/.config/zettelkasten/config.toml`, use:

```bash
zettelkasten init
```

Alternatively, you can use `-f <path>` to set where you want to store config
file. Remember that you will have to keep typing `-f` forever then.

## Your first note

To create new note with appropriate filename and header, in a directory defined
by configuration, use:

```bash
nvim $(zettelkasten new)
```

You are encouraged to create an alias for this or create new mapping in (n)vim
itself.

# Philosophy

- Note is a record of a thought, plan or goal. It may describe something you
  have read, something you have seen or somewhere you have been to.
- Note is written the way it _could_ be published.
- You do not collect notes for just for data hoarding. You revisit, distill or
  even remove notes.
- You do not strive for hierarchy. You link notes.
- To label notes, you use tags.
- To create collections, you use indices. E.g. "books", "travel", "ideas".
- Use plain text. It is searchable and it will last.
- Your notes are **you**. Do not let corporations read your mind. Do not store
  them anywhere outside your PC unless you can safely encrypt them. Do not trust
  any online services.
- Backup.

# (n)vim plugin

You can use any text editor, as long it can easily cooperate with `zettelkasten`
binary, i.e. you are able to define key mappings to create and search notes.
Just because I use `nvim`, I prepared dedicated plugin. To install it, use your
plugin manager, e.g. [`vim-plug`](https://github.com/junegunn/vim-plug):

```vim
Plug 'radiand/zettelkasten'
```

Currently it defines following mappings:

- `<leader>zk` to create new note (supports tab-completion of workspaces),
- `<leader>zf` to find notes (requires [fzf](https://github.com/junegunn/fzf)),
- `<leader>zg` to follow links (use when UID is under cursor).

