package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	optionMarkdown    = flag.Bool("markdown", false, "output shrink markdown")
	optionVerbose     = flag.Bool("v", false, "verbose")
	optionStyleFile   = flag.String("sf", "", "Specify a stylesheet file `path`")
	optionStyleInline = flag.String("st", "", "Specify `stylesheet` text directly as a string.")
	optionRootDir     = flag.String("d", ".", "Output `directory`")
)

func expandWildcard(args []string) []string {
	_args := []string{}
	for _, arg := range args {
		if matches, err := filepath.Glob(arg); err == nil && len(matches) >= 1 {
			_args = append(_args, matches...)
		} else {
			_args = append(_args, arg)
		}
	}
	return _args
}

func mains(args []string) error {
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

	if len(args) <= 0 {
		source, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if *optionMarkdown {
			if err := enexToMarkdown(*optionRootDir, "", source, "", verbose); err != nil {
				return err
			}
		} else {
			if err := enexToHtml(*optionRootDir, "", source, "", verbose); err != nil {
				return err
			}
		}
		return nil
	}
	_args := expandWildcard(args)

	if *optionMarkdown {
		return enexesToReadmeAndMarkdowns(*optionRootDir, verbose, _args)
	}
	return enexesToIndexAndHtmls(*optionRootDir, styleSheet, verbose, _args)
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
