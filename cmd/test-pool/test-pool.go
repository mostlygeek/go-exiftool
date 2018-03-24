package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	exiftool "github.com/mostlygeek/go-exiftool"
	"github.com/stretchr/powerwalk"
)

func main() {

	flag.Parse()
	root := flag.Arg(0)
	parallelism := runtime.NumCPU() * 2

	exif := exiftool.NewPool("exiftool", parallelism)
	//exif := exiftool.NewStayopen("exiftool")

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		metadata, err := exif.Extract(path)
		if err != nil {
			fmt.Println("FAIL", path, err.Error())
		} else {
			fmt.Println("OK", path, "mime:", metadata.MIMEType())
		}

		return nil
	}

	powerwalk.WalkLimit(root, walkFn, parallelism)
	exif.Stop()
}
