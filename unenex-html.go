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

func ToHtmls(rootDir, enexName string, source []byte, styleSheet string, wDebug, wLog io.Writer) error {
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

	opt := &Option{ExHeader: styleSheet}

	for _, note := range exports {
		html, imgSrc := note.Extract(opt)
		safeName := imgSrc.BaseName

		fmt.Fprintf(index, "<li><a href=\"%s\">%s</a></li>\n",
			url.PathEscape(safeName+".html"),
			note.Title,
		)
		fname := filepath.Join(rootDir, safeName+".html")
		fd, err := os.Create(fname)
		if err != nil {
			return err
		}
		io.WriteString(fd, html)
		fd.Close()
		fmt.Fprintln(wLog, "Create File:", fname)

		if err := extractAttachment(rootDir, imgSrc.Images, wLog); err != nil {
			return err
		}
	}
	return nil
}

func FilesToHtmls(rootDir, styleSheet string, enexFiles []string, wDebug, wLog io.Writer) error {
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
		if err := ToHtmls(rootDir, enexName, data, styleSheet, wDebug, wLog); err != nil {
			return err
		}
		fmt.Fprintf(wIndex, "<li><a href=\"%s/index.html\">%s</a></li>\n",
			url.PathEscape(enexName), enexName)
	}
	return nil
}
