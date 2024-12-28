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

func TestFindWidth(t *testing.T) {
	test := []struct {
		source string
		width  string
		height string
	}{
		{
			source: ` style="--en-naturalWidth:796; --en-naturalHeight:559;" hash="f3a35235096d45a300979dfee31ecda3" type="image/png" /`,
			width:  "796",
			height: "559",
		},
		{
			source: ` style="--en-naturalWidth:1920; --en-naturalHeight:1280;" height="384px" width="576px" hash="15426904081f2cfc80894e46d4e84723" type="image/jpeg" /`,
			width:  "576px",
			height: "384px",
		},
	}

	for _, c := range test {
		width := findWidth(c.source)
		if width != c.width {
			t.Fatalf(`expect "%s" as width for "%s", but "%s"`,
				c.width, c.source, width)
		}
		height := findHeight(c.source)
		if height != c.height {
			t.Fatalf(`expect "%s" as width for "%s", but "%s"`,
				c.height, c.source, height)
		}
	}
}
