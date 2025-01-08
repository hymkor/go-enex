package enex

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func makeDir(rootDir, name string, log io.Writer) error {
	if name == "" {
		return nil
	}
	name = filepath.Join(rootDir, name)
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

func getEnexBaseName(enexFileName string) string {
	enexName := filepath.Base(enexFileName)
	enexName = enexName[:len(enexName)-len(filepath.Ext(enexName))]
	// enexName =~ s/\.+$//;
	for len(enexName) > 0 && enexName[len(enexName)-1] == '.' {
		enexName = enexName[:len(enexName)-1]
	}
	return enexName
}
