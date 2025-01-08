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

func getEnexBaseName(enexFileName string) string {
	enexName := filepath.Base(enexFileName)
	enexName = enexName[:len(enexName)-len(filepath.Ext(enexName))]
	// enexName =~ s/\.+$//;
	for len(enexName) > 0 && enexName[len(enexName)-1] == '.' {
		enexName = enexName[:len(enexName)-1]
	}
	return enexName
}
