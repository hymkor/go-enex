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
	optionRootDir     = flag.String("d", ".", "Output directory")
)

func makeDir(root, name string, log io.Writer) error {
	if name == "" {
		return nil
	}
	name = filepath.Join(root, name)
	if stat, err := os.Stat(name); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		fmt.Fprintln(log, "Create Dir:", name)
		if err := os.Mkdir(name, 0755); err != nil {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("%s: fail to mkdir (file exists)", name)
	}
	return nil
}

func extractAttachment(root string, attachment map[string]*enex.Resource, log io.Writer) error {
	for _fname, data := range attachment {
		fname := filepath.Join(root, _fname)
		dir := filepath.Dir(fname)
		if stat, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Fprintln(log, "Create Dir:", dir)
			if err := os.Mkdir(dir, 0755); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if !stat.IsDir() {
			return fmt.Errorf("Can not mkdir %s because the file same name already exists", dir)
		}
		fmt.Fprint(log, "Create File: ", fname)
		fd, err := os.Create(fname)
		if err != nil {
			fmt.Fprintln(log)
			return fmt.Errorf("os.Create: %w", err)
		}
		n, err := data.WriteTo(fd)
		if err != nil {
			fmt.Fprintln(log)
			return fmt.Errorf(".WriteTo: %w", err)
		}
		if err := fd.Close(); err != nil {
			fmt.Fprintln(log)
			return fmt.Errorf(".Close: %w", err)
		}
		fmt.Fprintf(log, " (%d bytes)\n", n)
	}
	return nil
}

func enexToMarkdown(root, enexName string, source []byte, styleSheet string, verbose io.Writer) error {
	exports, err := enex.ParseMulti(source, verbose)
	if err != nil {
		return err
	}
	err = makeDir(root, enexName, os.Stderr)
	if err != nil {
		return err
	}
	root = filepath.Join(root, enexName)

	index, err := os.Create(filepath.Join(root, "README.md"))
	if err != nil {
		return err
	}
	defer index.Close()

	if enexName != "" {
		fmt.Fprintf(index, "# %s\n\n", enexName)
	}

	for _, note := range exports {
		safeName := enex.ToSafe.Replace(note.Title)

		fmt.Fprintf(index, "* [%s](%s)\n",
			note.Title,
			url.PathEscape(safeName+".md"),
		)

		html, imgSrc := note.HtmlAndDir()

		var markdown strings.Builder
		godown.Convert(&markdown, strings.NewReader(html), nil)
		fname := filepath.Join(root, safeName+".md")
		fd, err := os.Create(fname)
		if err != nil {
			return err
		}
		enex.ShrinkMarkdown(strings.NewReader(markdown.String()), fd)
		fd.Close()
		fmt.Fprintln(os.Stderr, "Create File:", fname)

		if err := extractAttachment(root, imgSrc.Images, os.Stderr); err != nil {
			return err
		}
	}
	return nil
}

const indexHtmlHeader = `<html><head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
</head><body>`

const indexHtmlFooter = "</body></html>"

func enexToHtml(root, enexName string, source []byte, styleSheet string, verbose io.Writer) error {
	exports, err := enex.ParseMulti(source, verbose)
	if err != nil {
		return err
	}
	err = makeDir(root, enexName, os.Stderr)
	if err != nil {
		return err
	}
	root = filepath.Join(root, enexName)

	index, err := os.Create(filepath.Join(root, "index.html"))
	if err != nil {
		return err
	}
	fmt.Fprintln(index, indexHtmlHeader)
	if enexName != "" {
		fmt.Fprintf(index, "<h1>%s</h1>\n\n", enexName)
	}
	fmt.Fprintln(index, "<ul>")
	defer func() {
		fmt.Fprintln(index, "</ul>")
		fmt.Fprintln(index, indexHtmlFooter)
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
		fname := filepath.Join(root, safeName+".html")
		fd, err := os.Create(fname)
		if err != nil {
			return err
		}
		io.WriteString(fd, html)
		fd.Close()
		fmt.Println("Create File:", fname)

		if err := extractAttachment(root, imgSrc.Images, os.Stderr); err != nil {
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
		if err := outfunc(*optionRootDir, "", source, "", verbose); err != nil {
			return err
		}
		return nil
	}
	_args := []string{}
	for _, arg := range args {
		if matches, err := filepath.Glob(arg); err == nil && len(matches) >= 1 {
			_args = append(_args, matches...)
		} else {
			_args = append(_args, arg)
		}
	}
	var fd *os.File
	if *optionMarkdown {
		fd, err = os.Create(filepath.Join(*optionRootDir, "README.md"))
		if err != nil {
			return err
		}
		defer fd.Close()
	} else {
		fd, err = os.Create(filepath.Join(*optionRootDir, "index.html"))
		if err != nil {
			return err
		}
		fmt.Fprintln(fd, indexHtmlHeader)
		fmt.Fprintln(fd, "<ul>")
		defer func() {
			fmt.Fprintln(fd, "</ul>")
			fmt.Fprintln(fd, indexHtmlFooter)
			fd.Close()
		}()
	}
	for _, arg := range _args {
		data, err = os.ReadFile(arg)
		if err != nil {
			return err
		}
		enexName := filepath.Base(arg)
		enexName = enexName[:len(enexName)-len(filepath.Ext(enexName))]
		if err := outfunc(*optionRootDir, enexName, data, styleSheet, verbose); err != nil {
			return err
		}
		if *optionMarkdown {
			fmt.Fprintf(fd, "- [%s](%s/README.md)\n",
				enexName, url.PathEscape(enexName))
		} else {
			fmt.Fprintf(fd, "<li><a href=\"%s/index.html\">%s</a></li>\n",
				url.PathEscape(enexName), enexName)
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
