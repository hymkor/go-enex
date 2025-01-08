package enex

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// ToMarkdowns converts an ENEX file to markdown format.
// The output markdown is saved under the directory specified by rootDir, with the file named enexName.
// The source contains the ENEX data, and the htmlToMarkdown function is used to convert HTML to markdown.
// Debug and log information are written to wDebug and wLog, respectively.
func ToMarkdowns(rootDir, enexName string, source []byte, htmlToMarkdown func(io.Writer, io.Reader) error, wDebug, wLog io.Writer) error {
	exports, err := Parse(source, wDebug)
	if err != nil {
		return err
	}
	err = makeDir(rootDir, enexName, wLog)
	if err != nil {
		return err
	}
	rootDir = filepath.Join(rootDir, enexName)

	index, err := os.Create(filepath.Join(rootDir, "README.md"))
	if err != nil {
		return err
	}
	defer index.Close()

	if enexName != "" {
		fmt.Fprintf(index, "# %s\n\n", enexName)
	}

	for _, note := range exports {
		html, bundle := note.Extract(nil)
		safeName := bundle.BaseName

		fmt.Fprintf(index, "* [%s](%s)\n",
			note.Title,
			url.PathEscape(safeName+".md"),
		)

		var markdown strings.Builder
		htmlToMarkdown(&markdown, strings.NewReader(html))
		fname := filepath.Join(rootDir, safeName+".md")
		fd, err := os.Create(fname)
		if err != nil {
			return err
		}
		shrinkMarkdown(strings.NewReader(markdown.String()), fd)
		fd.Close()
		fmt.Fprintln(wLog, "Create File:", fname)

		if err := bundle.Extract(rootDir, wLog); err != nil {
			return err
		}
	}
	return nil
}

// FilesToMarkdowns processes multiple enex files, converts each to markdown,
// and saves the output in the specified rootDir.
// The htmlToMarkdown function is used to convert HTML content to markdown.
// Debug and log information are written to wDebug and wLog, respectively.
func FilesToMarkdowns(rootDir string, htmlToMarkdown func(io.Writer, io.Reader) error, enexFiles []string, wDebug, wLog io.Writer) error {
	wReadme, err := os.Create(filepath.Join(rootDir, "README.md"))
	if err != nil {
		return err
	}
	defer wReadme.Close()

	for _, enexFileName := range enexFiles {
		data, err := os.ReadFile(enexFileName)
		if err != nil {
			return err
		}
		enexName := getEnexBaseName(enexFileName)
		if err := ToMarkdowns(rootDir, enexName, data, htmlToMarkdown, wDebug, wLog); err != nil {
			return err
		}
		fmt.Fprintf(wReadme, "- [%s](%s/README.md)\n",
			enexName, url.PathEscape(enexName))
	}
	return nil
}
