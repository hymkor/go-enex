[![GoDoc](https://godoc.org/github.com/hymkor/go-enex?status.svg)](https://godoc.org/github.com/hymkor/go-enex)

go-enex - Convert Evernote's export file(\*.enex) into HTML and images
==================================================

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
    export, err := enex.Parse(data)
    if err != nil {
        return err
    }
    html, images := export.Html("images-")
    fmt.Println(html)

    for fname, data := range images {
        fmt.Fprintf(os.Stderr, "Create File: %s (%d bytes)\n", fname, len(data))
        os.WriteFile(fname, data, 0666)
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

```
$ go run example.go < sample.enex > sample.html
Create File: images-image_19.png (7232 bytes)
Create File: images-image_21.png (3633 bytes)
Create File: images-image_4.png (50815 bytes)
Create File: images-image_9.png (54726 bytes)
Create File: images-image_11.png (52430 bytes)
Create File: images-image_13.png (52293 bytes)
```

Executable
-----------

### Install

Download the binary package from [Releases](https://github.com/hymkor/go-enex/releases) and extract the executable.

#### for scoop-installer

```
scoop install https://raw.githubusercontent.com/hymkor/go-enex/master/enexToHtml.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install enexToHtml
```

#### Example

```
$ cd cmd/enexToHtml
$ go build
```

```
$ ./enexToHtml           [-prefix=PREFIX] < ENEX-FILENAME.enex > ENEX-FILENAME.html
$ ./enexToHtml -markdown [-prefix=PREFIX] < ENEX-FILENAME.enex > ENEX-FILENAME.md
$ ./enexToHtml [-markdown] [-prefix=PREFIX] ENEX-FILENAME.enex
```

- The PREFIX is used as filename-header for image-files.
- `-markdown` option is by [mattn/go-godown](https://github.com/mattn/godown)
