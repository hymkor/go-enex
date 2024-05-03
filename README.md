[![GoDoc](https://godoc.org/github.com/hymkor/go-enex?status.svg)](https://godoc.org/github.com/hymkor/go-enex)

go-enex - Convert Evernote's export file(\*.enex) into HTML and images
==================================================

How to use executable (unenex)
-------------------------

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
$ ./unenex [-markdown] ENEX-FILENAME.enex
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
    // enex.Parse can not support multi-notes enex file.
    // To Parse multi-notes enex file, use `enex.ParseMulti`
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
