//go:build ignore

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

		if len(bundle.Images) > 0 {
			fmt.Fprintf(os.Stderr, "Create Dir: %s", bundle.Dir)
			os.Mkdir(bundle.Dir, 0755)
			for fname, rsc := range bundle.Images {
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
