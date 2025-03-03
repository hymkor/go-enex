package enex

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

const indexHtmlHeader = `<!DOCTYPE html><html><head>
<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
</head><body>`

const indexHtmlFooter = "</body></html>"

// ToHtmls converts a single ENEX file into HTML format.
// The output HTML is saved under the directory specified by rootDir, with enexName as the note name.
// The content is read from source, the styleSheet is applied to the HTML,
// and debug and log information is written to wDebug and wLog, respectively.
// If webClipOnly is true, only the web-clip content will be output without Evernote styling.
func ToHtmls(rootDir, enexName string, source []byte, styleSheet string, webClipOnly bool, wDebug, wLog io.Writer) error {
	exports, err := Parse(source, wDebug)
	if err != nil {
		return err
	}
	err = makeDir(rootDir, enexName, wLog)
	if err != nil {
		return err
	}
	rootDir = filepath.Join(rootDir, enexName)

	index, err := os.Create(filepath.Join(rootDir, "index.html"))
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

	opt := &Option{
		ExHeader:    styleSheet,
		WebClipOnly: webClipOnly,
	}

	for _, note := range exports {
		html, bundle := note.Extract(opt)
		safeName := bundle.BaseName

		fmt.Fprintf(index, "<li><a href=\"%s\">%s</a></li>\n",
			url.PathEscape(safeName+".html"),
			note.Title,
		)
		fname := filepath.Join(rootDir, safeName+".html")
		if err := os.WriteFile(fname, []byte(html), 0644); err != nil {
			return err
		}
		fmt.Fprintln(wLog, "Create File:", fname)

		if err := bundle.Extract(rootDir, wLog); err != nil {
			return err
		}
	}
	return nil
}

// FilesToHtmls converts multiple ENEX files into HTML format.
// The output HTML files are saved under the directory specified by rootDir, with each ENEX file being processed.
// The styleSheet is applied to the HTML, and debug and log information are written to wDebug and wLog, respectively.
// If webClipOnly is true, only the web-clip content will be output without Evernote styling.
func FilesToHtmls(rootDir, styleSheet string, enexFiles []string, webClipOnly bool, wDebug, wLog io.Writer) error {
	wIndex, err := os.Create(filepath.Join(rootDir, "index.html"))
	if err != nil {
		return err
	}
	fmt.Fprintln(wIndex, indexHtmlHeader)
	fmt.Fprintln(wIndex, "<ul>")
	defer func() {
		fmt.Fprintln(wIndex, "</ul>")
		fmt.Fprintln(wIndex, indexHtmlFooter)
		defer wIndex.Close()
	}()

	for _, enexFileName := range enexFiles {
		data, err := os.ReadFile(enexFileName)
		if err != nil {
			return err
		}
		enexName := getEnexBaseName(enexFileName)
		if err := ToHtmls(rootDir, enexName, data, styleSheet, webClipOnly, wDebug, wLog); err != nil {
			return err
		}
		fmt.Fprintf(wIndex, "<li><a href=\"%s/index.html\">%s</a></li>\n",
			url.PathEscape(enexName), enexName)
	}
	return nil
}
