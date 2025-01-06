[![GoDoc](https://godoc.org/github.com/hymkor/go-enex?status.svg)](https://pkg.go.dev/github.com/hymkor/go-enex)

Go-enex &amp; Unenex
====================

Convert Evernote's export file(\*.enex) into HTML(or markdown) and images.

- `go-enex` : the package for Go
- `unenex` : the executable using go-enex

How to use Unenex
---------------------

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
$ unenex [-markdown] ENEX-FILENAME.enex
```

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
        html, imgSrc := note.Extract(nil)
        baseName := imgSrc.BaseName
        err := os.WriteFile(baseName+".html", []byte(html), 0644)
        if err != nil {
            return err
        }
        fmt.Fprintf(os.Stderr, "Create File: %s.html (%d bytes)\n", baseName, len(html))

        if len(imgSrc.Images) > 0 {
            fmt.Fprintf(os.Stderr, "Create Dir: %s", imgSrc.Dir)
            os.Mkdir(imgSrc.Dir, 0755)
            for fname, rsc := range imgSrc.Images {
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

- [ReleaseNote(en)](release_note_en.md)
- [ReleaseNote(ja)](release_note_ja.md)
