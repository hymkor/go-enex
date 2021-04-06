package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/go-enex"
)

func mains() error {
	en, err := enex.Parse(os.Stdin)
	if err != nil {
		return err
	}
	html, attachment := en.Html("")
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
