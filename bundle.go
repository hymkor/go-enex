package enex

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type Attachments struct {
	Images    map[string]*Resource
	BaseName  string
	Dir       string
	dirEscape string
	serialNo  _SerialNo
	sanitizer func(string) string
}

func newAttachments(note *Note, sanitizer func(string) string) *Attachments {
	baseName := sanitizer(note.Title)
	dir := baseName + ".files"
	dirEscape := url.PathEscape(dir)
	return &Attachments{
		Images:    make(map[string]*Resource),
		BaseName:  baseName,
		Dir:       dir,
		dirEscape: dirEscape,
		serialNo:  make(_SerialNo),
		sanitizer: sanitizer,
	}
}

func (attach *Attachments) makeUrlFor(rsc *Resource) string {
	name := attach.sanitizer(attach.serialNo.ToUniqName(rsc.Mime, rsc.FileName, rsc.Hash))
	rsc.NewFileName = name
	attach.Images[filepath.Join(attach.Dir, name)] = rsc
	return path.Join(attach.dirEscape, url.PathEscape(name))
}

func (A *Attachments) Extract(rootDir string, log io.Writer) error {
	attachment := A.Images
	for _fname, data := range attachment {
		fname := filepath.Join(rootDir, _fname)
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
		data, err := data.Data()
		if err != nil {
			return err
		}
		if err := os.WriteFile(fname, data, 0644); err != nil {
			return err
		}
		fmt.Fprintf(log, "Create File: %s (%d bytes)\n", fname, len(data))
	}
	return nil
}
