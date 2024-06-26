package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mattn/godown"

	"github.com/hymkor/go-enex"
)

var optionMarkdown = flag.Bool("markdown", false, "output shrink markdown")

var optionVerbose = flag.Bool("v", false, "verbose")

var unmarkdown = strings.NewReplacer()

var toSafe = strings.NewReplacer(
	`<`, `＜`,
	`>`, `＞`,
	`"`, `”`,
	`/`, `／`,
	`\`, `＼`,
	`|`, `｜`,
	`?`, `？`,
	`*`, `＊`,
	`:`, `：`,
	`(`, `（`,
	`)`, `）`,
	` `, `_`,
	// `.`, `．`,
)

func mains(args []string) error {
	var data []byte
	var err error

	if len(args) <= 0 {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		data, err = os.ReadFile(args[0])
		if err != nil {
			return err
		}
	}
	verbose := io.Discard
	if *optionVerbose {
		verbose = os.Stderr
	}
	exports, err := enex.ParseMulti(data, verbose)
	if err != nil {
		return err
	}
	var index *os.File
	if *optionMarkdown {
		var err error
		index, err = os.Create("README.md")
		if err != nil {
			return err
		}
		defer index.Close()
	} else {
		var err error
		index, err = os.Create("index.html")
		if err != nil {
			return err
		}
		fmt.Fprintln(index, `<html lang="ja"><head>`)
		fmt.Fprintln(index, `<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">`)
		fmt.Fprintln(index, `</head><body><ul>`)
		defer func() {
			fmt.Fprintln(index, "</ul></body></html>")
			index.Close()
		}()
	}
	for _, note := range exports {
		var html string
		var images map[string]*enex.Resource

		safeName := toSafe.Replace(note.Title)

		if *optionMarkdown {
			fmt.Fprintf(index, "* [%s](%s)\n",
				note.Title,
				url.QueryEscape(safeName+".md"),
			)
		} else {
			fmt.Fprintf(index, "<li><a href=\"%s\">%s</a></li>\n",
				url.QueryEscape(safeName+".html"),
				note.Title,
			)
		}
		html, images = note.HtmlAndImagesWithRenamer(
			func(fname string, index int) string {
				ext := filepath.Ext(fname)
				baseName := toSafe.Replace(fname[:len(fname)-len(ext)])
				fname = fmt.Sprintf("%s_%d%s", baseName, index, ext)
				fullpath := path.Join(safeName+".files", fname)
				return fullpath
			},
		)
		if *optionMarkdown {
			var markdown strings.Builder
			godown.Convert(&markdown, strings.NewReader(html), nil)
			fd, err := os.Create(safeName + ".md")
			if err != nil {
				return err
			}
			enex.ShrinkMarkdown(strings.NewReader(markdown.String()), fd)
			fd.Close()
			fmt.Println("Create File:", safeName+".md")
		} else {
			fd, err := os.Create(safeName + ".html")
			if err != nil {
				return err
			}
			io.WriteString(fd, html)
			fd.Close()
			fmt.Println("Create File:", safeName+".html")
		}
		for fname, data := range images {
			dir := filepath.Dir(fname)
			if stat, err := os.Stat(dir); os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, "Create Dir", dir)
				if err := os.Mkdir(dir, 0755); err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else if !stat.IsDir() {
				return fmt.Errorf("Can not mkdir %s because the file same name already exists", dir)
			}
			fmt.Fprint(os.Stderr, "Create File:", fname)
			fd, err := os.Create(fname)
			if err != nil {
				fmt.Fprintln(os.Stderr)
				return fmt.Errorf("os.Create: %w", err)
			}
			n, err := data.WriteTo(fd)
			if err != nil {
				fmt.Fprintln(os.Stderr)
				return fmt.Errorf(".WriteTo: %w", err)
			}
			if err := fd.Close(); err != nil {
				fmt.Fprintln(os.Stderr)
				return fmt.Errorf(".Close: %w", err)
			}
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
