go-enex - evernote's export data to HTML converter
==================================================

Library for Go
--------------

```go
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
