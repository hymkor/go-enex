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

var optionShrink = flag.Bool("shrink-markdown", false, "output shrink markdown")

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
		if *optionMarkdown || *optionShrink {
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
	html, images := export.HtmlAndImagesWithRenamer(enex.DefaultRenamer(baseName))
	if *optionShrink {
		var markdown strings.Builder
		godown.Convert(&markdown, strings.NewReader(html), nil)
		enex.ShrinkMarkdown(strings.NewReader(markdown.String()), output)
	} else if *optionMarkdown {
		godown.Convert(output, strings.NewReader(html), nil)
	} else {
		io.WriteString(output, html)
	}
	for fname, data := range images {
		fmt.Fprint(os.Stderr, "Create File: ", fname)
		fd, err := os.Create(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr)
			return err
		}
		n, _ := data.WriteTo(fd)
		fd.Close()
		fmt.Fprintf(os.Stderr, " (%d bytes)\n", n)
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
