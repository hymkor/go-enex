package enex

import (
	"testing"
)

func TestSerialNoToUniqName(t *testing.T) {
	S := _SerialNo{}

	source := []struct {
		name string
		hash string
	}{
		{"hogehoge.txt", "0"},
		{"uhauha.txt", "1"},
		{"hogehoge.txt", "2"},
	}

	expect := []string{
		"hogehoge.txt",
		"uhauha.txt",
		"hogehoge (1).txt",
	}

	for i, src := range source {
		result := S.ToUniqName("text/plain", src.name, src.hash)
		if expect[i] != result {
			t.Fatalf(`(%d) expect "%s" for "%s" (%s), but "%s"`,
				i, expect[i], src.name, src.hash, result)
		}
	}
}

func TestFindWidth(t *testing.T) {
	test := []struct {
		source string
		width  string
		height string
	}{
		{
			source: ` height="384px" width="576px" hash="15426904081f2cfc80894e46d4e84723" type="image/jpeg" /`,
			width:  "576px",
			height: "384px",
		},
	}

	for _, c := range test {
		a := parseEnMediaAttr(c.source)
		width := a["width"]
		if width != c.width {
			t.Fatalf(`expect "%s" as width for "%s", but "%s"`,
				c.width, c.source, width)
		}
		height := a["height"]
		if height != c.height {
			t.Fatalf(`expect "%s" as width for "%s", but "%s"`,
				c.height, c.source, height)
		}
	}
}
