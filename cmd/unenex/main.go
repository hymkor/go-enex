package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mattn/godown"

	"github.com/hymkor/go-enex"
)

var (
	optionMarkdown    = flag.Bool("markdown", false, "output shrink markdown")
	optionVerbose     = flag.Bool("v", false, "verbose")
	optionStyleFile   = flag.String("sf", "", "Specify a stylesheet file")
	optionStyleInline = flag.String("st", "", "Specify stylesheet text directly as a string.")
)

func makeAndChdir(name string) (func(), error) {
	if name == "" {
		return func() {}, nil
	}
	if _, err := os.Stat(name); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Create Dir:", name)
		if err := os.Mkdir(name, 0755); err != nil {
			return nil, err
		}
	}
	if err := os.Chdir(name); err != nil {
		return nil, err
	}
	return func() { os.Chdir("..") }, nil
}

func extractAttachment(attachment map[string]*enex.Resource) error {
	for fname, data := range attachment {
		dir := filepath.Dir(fname)
		if stat, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Create Dir:", dir)
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if !stat.IsDir() {
			return fmt.Errorf("Can not mkdir %s because the file same name already exists", dir)
		}
		fmt.Fprint(os.Stderr, "Create File: ", fname)
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
	return nil
}

func enexToMarkdown(name string, source []byte, styleSheet string, verbose io.Writer) error {
	exports, err := enex.ParseMulti(source, verbose)
	if err != nil {
		return err
	}
	closer, err := makeAndChdir(name)
	if err != nil {
		return err
	}
	defer closer()

	index, err := os.Create("README.md")
	if err != nil {
		return err
	}
	defer index.Close()

	for _, note := range exports {
		safeName := enex.ToSafe.Replace(note.Title)

		fmt.Fprintf(index, "* [%s](%s)\n",
			note.Title,
			url.PathEscape(safeName+".md"),
		)

		html, imgSrc := note.HtmlAndDir()

		var markdown strings.Builder
		godown.Convert(&markdown, strings.NewReader(html), nil)
		fd, err := os.Create(safeName + ".md")
		if err != nil {
			return err
		}
		enex.ShrinkMarkdown(strings.NewReader(markdown.String()), fd)
		fd.Close()
		fmt.Println("Create File:", safeName+".md")

		if err := extractAttachment(imgSrc.Images); err != nil {
			return err
		}
	}
	return nil
}

func enexToHtml(name string, source []byte, styleSheet string, verbose io.Writer) error {
	exports, err := enex.ParseMulti(source, verbose)
	if err != nil {
		return err
	}
	closer, err := makeAndChdir(name)
	if err != nil {
		return err
	}
	defer closer()

	index, err := os.Create("index.html")
	if err != nil {
		return err
	}
	fmt.Fprintln(index, `<html><head>`)
	fmt.Fprintln(index, `<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">`)
	fmt.Fprintln(index, `</head><body><ul>`)
	defer func() {
		fmt.Fprintln(index, "</ul></body></html>")
		index.Close()
	}()

	for _, note := range exports {
		note.ExHeader = styleSheet

		safeName := enex.ToSafe.Replace(note.Title)

		fmt.Fprintf(index, "<li><a href=\"%s\">%s</a></li>\n",
			url.PathEscape(safeName+".html"),
			note.Title,
		)
		html, imgSrc := note.HtmlAndDir()
		fd, err := os.Create(safeName + ".html")
		if err != nil {
			return err
		}
		io.WriteString(fd, html)
		fd.Close()
		fmt.Println("Create File:", safeName+".html")

		if err := extractAttachment(imgSrc.Images); err != nil {
			return err
		}
	}
	return nil
}

func mains(args []string) error {
	var data []byte
	var err error

	verbose := io.Discard
	if *optionVerbose {
		verbose = os.Stderr
	}
	var styleSheet string
	if *optionStyleFile != "" {
		var buffer strings.Builder
		fd, err := os.Open(*optionStyleFile)
		if err != nil {
			return err
		}
		buffer.WriteString("<style>\n")
		io.Copy(&buffer, fd)
		if *optionStyleInline != "" {
			fmt.Fprintln(&buffer)
			fmt.Fprintln(&buffer, *optionStyleInline)
		}
		buffer.WriteString("\n</style>\n")
		styleSheet = buffer.String()
		fd.Close()
	} else if *optionStyleInline != "" {
		styleSheet = fmt.Sprintf("<style>\n%s\n</style>\n", *optionStyleInline)
	}

	outfunc := enexToHtml
	if *optionMarkdown {
		outfunc = enexToMarkdown
	}

	if len(args) <= 0 {
		source, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if err := outfunc("", source, "", verbose); err != nil {
			return err
		}
	} else {
		_args := []string{}
		for _, arg := range args {
			if matches, err := filepath.Glob(arg); err == nil && len(matches) >= 1 {
				_args = append(_args, matches...)
			} else {
				_args = append(_args, arg)
			}
		}
		for _, arg := range _args {
			data, err = os.ReadFile(arg)
			if err != nil {
				return err
			}
			enexName := filepath.Base(arg)
			enexName = enexName[:len(enexName)-len(filepath.Ext(enexName))]
			if err := outfunc(enexName, data, styleSheet, verbose); err != nil {
				return err
			}
		}
	}
	return nil
}

var version string

func main() {
	fmt.Fprintf(flag.CommandLine.Output(), "%s %s-%s-%s\n",
		filepath.Base(os.Args[0]),
		version,
		runtime.GOOS,
		runtime.GOARCH)
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
