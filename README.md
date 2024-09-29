_zettelkasten_ – creating and managing plain text notes.

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
are kept as .md files, but you can just drop few sentences without further
formatting. As long as your editor of choice supports it, dropping here an
![image](./resources/image.jpg) might be possible.
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
go install .
```

If your OS is setup properly and your `PATH` is populated with `go` built
binaries, you can use `zettelkasten` in your terminal now.

## Initialize

To follow default settings and create configuration in
`$HOME/.config/zettelkasten/config.toml`, use:

```bash
zettelkasten init
```

Alternatively, you can use `-config <path>` to set where you want to store
config file. Remember that you will have to keep typing `-config` forever then.

## Your first note

To create new note with appropriate filename and header, in a directory defined
by configuration, use:

```bash
nvim $(zettelkasten new -stdout=false)
```

You are encouraged to create an alias for this or create new mapping in (n)vim
itself.
