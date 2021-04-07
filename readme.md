go-enex - evernote's export data to HTML converter
==================================================

Library for Go
--------------

```go
// +build ignore

package main

import (
    "fmt"
    "io"
    "os"

    "github.com/zetamatta/go-enex"
)

func mains() error {
    data, err := io.ReadAll(os.Stdin)
    if err != nil {
        return err
    }
    en, err := enex.Parse(data)
    if err != nil {
        return err
    }
    html, attachment := en.Html("attachment-")
    io.WriteString(os.Stdout, html)

    for fname, data := range attachment {
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
Create File: attachment-image_1.png (21957 bytes)
Create File: attachment-image_7.png (17981 bytes)
Create File: attachment-image_14.png (56557 bytes)
Create File: attachment-image_13.png (52293 bytes)
Create File: attachment-image.png (37202 bytes)
Create File: attachment-image_5.png (27353 bytes)
Create File: attachment-image_6.png (57539 bytes)
```

Executable
-----------

```
$ cd cmd/enexToHtml
$ go build
```

```
$ ./enexToHtml           [-prefix=PREFIX] < ENEX-FILENAME.enex > ENEX-FILENAME.html
$ ./enexToHtml -markdown [-prefix=PREFIX] < ENEX-FILENAME.enex > ENEX-FILENAME.md
$ ./enexToHtml [-markdown] [-prefix=PREFIX] ENEX-FILENAME.enex
```

The PREFIX is used as filename-header for attachment-image-files.
