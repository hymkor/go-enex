package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zetamatta/go-enex"
)

func mains(args []string) error {
	var data []byte
	var err error
	var output io.Writer
	prefix := ""

	if len(args) <= 0 {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		output = os.Stdout
	} else {
		data, err = os.ReadFile(args[0])
		if err != nil {
			return err
		}
		ext := filepath.Ext(args[0])
		prefix = args[0][:len(args[0])-len(ext)]
		fd, err := os.Create(prefix + ".html")
		if err != nil {
			return err
		}
		defer fd.Close()
		output = fd
		prefix = prefix + "-"
	}
	en, err := enex.Parse(data)
	if err != nil {
		return err
	}
	html, attachment := en.Html(prefix)
	io.WriteString(output, html)
	for fname, data := range attachment {
		fmt.Fprintf(os.Stderr, "Create File: %s (%d bytes)\n", fname, len(data))
		os.WriteFile(fname, data, 0666)
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
