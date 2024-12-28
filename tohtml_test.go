package enex

import (
	"testing"
)

func TestSerialNoToUniqName(t *testing.T) {
	S := SerialNo{}

	source := []struct {
		name  string
		index int
	}{
		{"hogehoge.txt", 0},
		{"uhauha.txt", 1},
		{"hogehoge.txt", 2},
	}

	expect := []string{
		"hogehoge.txt",
		"uhauha.txt",
		"hogehoge (2).txt",
	}

	for i, src := range source {
		result := S.ToUniqName(src.name, src.index)
		if expect[i] != result {
			t.Fatalf(`(%d) expect "%s" for "%s" #%d, but "%s"`,
				i, expect[i], src.name, src.index, result)
		}
	}
}
