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
    "net/url"
    "os"
    "path"
    "path/filepath"
    "strings"

    "github.com/hymkor/go-enex"
)

var toSafe = strings.NewReplacer(
    `<`, `＜`,
    `>`, `＞`,
    `"`, `”`,
    `/`, `／`,
    `\`, `＼`,
    `|`, `｜`,
    `?`, `？`,
    `*`, `＊`,
    `:`, `：`,
    `(`, `（`,
    `)`, `）`,
    ` `, `_`,
)

func toUniqName(name string, index int) string {
    ext := path.Ext(name)
    base := name[:len(name)-len(ext)]
    return fmt.Sprintf("%s%d%s", base, index, ext)
}

func mains() error {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        return err
    }
    notes, err := enex.ParseMulti(data, os.Stderr)
    if err != nil {
        return err
    }
    for _, note := range notes {
        baseName := toSafe.Replace(note.Title)
        images := make(map[string]*enex.Resource)
        dir := baseName + ".files"
        dirEscape := url.PathEscape(dir)
        html := note.ToHtml(func(rsc *enex.Resource) string {
            name := toSafe.Replace(toUniqName(rsc.FileName, rsc.Index))
            images[filepath.Join(dir, name)] = rsc
            return path.Join(dirEscape, url.PathEscape(name))
        })
        err := os.WriteFile(baseName+".html", []byte(html), 0644)
        if err != nil {
            return err
        }
        fmt.Fprintf(os.Stderr, "Create File: %s.html (%d bytes)\n", baseName, len(html))

        if len(images) > 0 {
            os.Mkdir(dir, 0755)
            for fname, rsc := range images {
                data := rsc.Data()
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
