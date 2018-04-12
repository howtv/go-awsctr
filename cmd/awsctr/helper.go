package main

import (
	"bytes"
	"flag"
)

func HelpMsg(b *bytes.Buffer, line string, fs *flag.FlagSet) string {
	b.WriteString(line + "\n")
	fs.SetOutput(b)
	fs.PrintDefaults()
	return b.String()
}
