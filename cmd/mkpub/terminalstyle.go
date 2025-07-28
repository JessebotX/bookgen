package main

import (
	"strings"
)

const (
	TerminalClear       = "\033[0m"
	TerminalTextBold    = "1"
	TerminalTextRed     = "31"
	TerminalTextGreen   = "32"
	TerminalTextYellow  = "33"
	TerminalTextBlue    = "34"
	TerminalTextMagenta = "35"
	TerminalTextCyan    = "36"
	TerminalTextWhite   = "37"
)

func TerminalStyle(s string, codes ...string) string {
	if len(codes) == 0 {
		return s
	}

	code := strings.Join(codes, ";")

	return "\033[" + code + "m" + s + TerminalClear
}
