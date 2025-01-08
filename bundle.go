package enex

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type Bundle struct {
	Resource  map[string]*Resource
	BaseName  string
	Dir       string
	dirEscape string
	serialNo  _SerialNo
	sanitizer func(string) string
}

func newBundle(note *Note, sanitizer func(string) string) *Bundle {
	baseName := sanitizer(note.Title)
	dir := baseName + ".files"
	dirEscape := url.PathEscape(dir)
	return &Bundle{
		Resource:  make(map[string]*Resource),
		BaseName:  baseName,
		Dir:       dir,
		dirEscape: dirEscape,
		serialNo:  make(_SerialNo),
		sanitizer: sanitizer,
	}
}

func (B *Bundle) makeUrlFor(rsc *Resource) string {
	name := B.sanitizer(B.serialNo.ToUniqName(rsc.Mime, rsc.FileName, rsc.Hash))
	rsc.NewFileName = name
	B.Resource[filepath.Join(B.Dir, name)] = rsc
	return path.Join(B.dirEscape, url.PathEscape(name))
}

func (B *Bundle) Extract(rootDir string, log io.Writer) error {
	for _fname, data := range B.Resource {
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
