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
