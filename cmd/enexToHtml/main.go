package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/godown"
	"github.com/zetamatta/go-enex"
)

var optionMarkdown = flag.Bool("markdown", false, "output markdown")

var optionPrefix = flag.String("prefix", "", "prefix for attachement")

func mains(args []string) error {
	var data []byte
	var err error
	var output io.Writer
	baseName := ""

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
		baseName = args[0][:len(args[0])-len(ext)]
		outputSuffix := ".html"
		if *optionMarkdown {
			outputSuffix = ".md"
		}
		fd, err := os.Create(baseName + outputSuffix)
		if err != nil {
			return err
		}
		defer fd.Close()
		output = fd
		if *optionPrefix != "" {
			baseName = *optionPrefix
		} else {
			baseName = baseName + "-"
		}
	}
	export, err := enex.Parse(data)
	if err != nil {
		return err
	}
	html, images := export.Html(baseName)
	if *optionMarkdown {
		var markdown strings.Builder
		godown.Convert(&markdown, strings.NewReader(html), &godown.Option{})
		enex.ShrinkMarkdown(strings.NewReader(markdown.String()), output)
	} else {
		io.WriteString(output, html)
	}
	for fname, data := range images {
		fmt.Fprintf(os.Stderr, "Create File: %s (%d bytes)\n", fname, len(data))
		os.WriteFile(fname, data, 0666)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
