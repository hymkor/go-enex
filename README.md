[![GoDoc](https://godoc.org/github.com/hymkor/go-enex?status.svg)](https://pkg.go.dev/github.com/hymkor/go-enex)

Go-enex &amp; Unenex
====================

Convert Evernote's export file(\*.enex) into HTML(or markdown) and images.

- `go-enex` : the package for Go
- `unenex` : the executable using go-enex

How to use Unenex
-----------------

### Install

Download the binary package from [Releases](https://github.com/hymkor/go-enex/releases) and extract the executable.

### Use `go install`

```
go install github.com/hymkor/go-enex/cmd/unenex@latest
```

#### Use scoop-installer

```
scoop install https://raw.githubusercontent.com/hymkor/go-enex/master/unenex.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install unenex
```

#### Example

```
$ unenex [-markdown] {ENEX-FILENAME.enex}
```

- Square brackets `[ ]` indicate optional arguments.
- Curly braces `{ }` indicate that the argument can be repeated.
- `-markdown` makes a makedown file instead of HTML

Library for Go
--------------

```example.go
package main

import (
    "fmt"
    "io"
    "os"

    "github.com/hymkor/go-enex"
)

func mains() error {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        return err
    }
    notes, err := enex.Parse(data, os.Stderr)
    if err != nil {
        return err
    }
    for _, note := range notes {
        html, bundle := note.Extract(nil)
        baseName := bundle.BaseName
        err := os.WriteFile(baseName+".html", []byte(html), 0644)
        if err != nil {
            return err
        }
        fmt.Fprintf(os.Stderr, "Create File: %s.html (%d bytes)\n", baseName, len(html))

        if len(bundle.Resource) > 0 {
            fmt.Fprintf(os.Stderr, "Create Dir: %s", bundle.Dir)
            os.Mkdir(bundle.Dir, 0755)
            for fname, rsc := range bundle.Resource {
                data, err := rsc.Data()
                if err != nil {
                    return err
                }
                fmt.Fprintf(os.Stderr, "Create File: %s (%d bytes)\n", fname, len(data))
                os.WriteFile(fname, data, 0666)
            }
        }
    }
    return nil
}
func main() {
    if err := mains(); err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }
}
```

License
-------

MIT License

Acknowledgements
---------------

- [@Laetgark](https://github.com/Laetgark)
- [@Juelicher-Trainee](https://github.com/Juelicher-Trainee)

Release notes
-------------

- [Release Notes (English)](release_note_en.md)
- [Release Notes (Japanese)](release_note_ja.md)
