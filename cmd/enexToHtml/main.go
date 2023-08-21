package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/godown"

	"github.com/hymkor/go-enex"
)

var optionMarkdown = flag.Bool("markdown", false, "output shrink markdown")

var optionPrefix = flag.String("prefix", "", "prefix for attachement")

var optionEmbed = flag.Bool("embed", false, "use <img src=\"data:...\">")

var optionVerbose = flag.Bool("v", false, "verbose")

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
	verbose := io.Discard
	if *optionVerbose {
		verbose = os.Stderr
	}
	export, err := enex.Parse(data, verbose)
	if err != nil {
		return err
	}
	var html string
	var images map[string]*enex.Resource

	if *optionEmbed {
		html, images = export.HtmlAndImagesWithRenamer(
			func(name string, index int) string {
				rsc := export.Resource[name][index]
				return fmt.Sprintf("data:%s;base64,%s",
					rsc.Mime,
					strings.TrimSpace(strings.ReplaceAll(rsc.DataBeforeDecoded(), "\n", "")))
			})
	} else {
		html, images = export.HtmlAndImagesWithRenamer(enex.DefaultRenamer(baseName))
	}
	if *optionMarkdown {
		var markdown strings.Builder
		godown.Convert(&markdown, strings.NewReader(html), nil)
		enex.ShrinkMarkdown(strings.NewReader(markdown.String()), output)
	} else {
		io.WriteString(output, html)
	}
	if !*optionEmbed {
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
