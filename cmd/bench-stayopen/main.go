package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	exiftool "github.com/mostlygeek/go-exiftool"
)

func main() {
	flag.Parse()

	dir := flag.Arg(0)
	if info, err := os.Lstat(dir); err != nil || !info.IsDir() {
		fmt.Println("Error: arg0 not a directory")
		os.Exit(1)
	}

	et, err := exiftool.NewStayOpen("exiftool")
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	defer et.Stop()

	fmt.Println("Starting.... ")
	start := time.Now()
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		et.Extract(path)
		return nil
	})

	if err != nil {
		fmt.Println("Walk ERROR: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Took: %v\n", time.Now().Sub(start))
}
