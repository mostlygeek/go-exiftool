package main

import (
	"flag"
	"fmt"
	"os"

	exiftool "github.com/mostlygeek/go-exiftool"
)

func main() {

	flag.Parse()
	meta, err := exiftool.Extract(flag.Arg(0))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		fmt.Println("Mimetype: ", meta.MIMEType())
		j, err := meta.MarshalJSON()
		if err != nil {
			fmt.Println("Failed getting JSON", err.Error())
		} else {
			fmt.Println(string(j))
		}
	}
}
