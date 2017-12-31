package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	exiftool "github.com/mostlygeek/go-exiftool"
)

func main() {
	flag.Parse()
	meta, err := exiftool.Extract(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	} else {
		thumb, err := meta.GetBytes("ThumbnailImage")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		} else {
			io.Copy(os.Stdout, bytes.NewBuffer(thumb))
		}
	}
}
