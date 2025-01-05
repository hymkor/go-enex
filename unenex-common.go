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

func extractAttachment(rootDir string, attachment map[string]*Resource, log io.Writer) error {
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

func getEnexBaseName(enexFileName string) string {
	enexName := filepath.Base(enexFileName)
	enexName = enexName[:len(enexName)-len(filepath.Ext(enexName))]
	// enexName =~ s/\.+$//;
	for len(enexName) > 0 && enexName[len(enexName)-1] == '.' {
		enexName = enexName[:len(enexName)-1]
	}
	return enexName
}
