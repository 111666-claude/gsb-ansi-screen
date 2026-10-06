// Command ansi 是屏幕解析的场景入口。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"example.com/ansi-screen/term"
)

func emit(v any) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	sample := flag.String("sample", "", "cursor|clamp|clear|work")
	flag.Parse()

	c := &term.Counter{}
	switch *sample {
	case "cursor":
		s := term.New(5, 1)
		s.Feed("a\x1b[1;3Hb", c)
		emit(map[string]any{"rows": s.Text()})
	case "clamp":
		s := term.New(3, 1)
		s.Feed("\x1b[1;9Hx", c)
		emit(map[string]any{"rows": s.Text()})
	case "clear":
		s := term.New(6, 1)
		s.Feed("abc\x1b[2Jx", c)
		emit(map[string]any{"rows": s.Text()})
	case "work":
		s := term.New(80, 24)
		s.Feed(strings.Repeat("x", 2000), c)
		emit(map[string]int{"cells": 80 * 24, "scanned": c.Scanned})
	default:
		fmt.Fprintln(os.Stderr, "未知场景："+*sample)
		os.Exit(2)
	}
}
