[![](https://godoc.org/github.com/mostlygeek/parhash?status.svg)](https://godoc.org/github.com/mostlygeek/go-exiftool)

# About

go-exiftool makes it easy to extract metadata with [exiftool](https://sno.phy.queensu.ca/~phil/exiftool/) and work with it in Go.  There are currently no comparable native Go libraries with the breadth and depth of exiftool. In exchange for functionality there is a performance and a deployment penalty. 

Fortunately, these are minimal. exiftool only requires perl5, which is available by default on almost every platform. The performance overhead of using an external program can be mitigated in many ways (ie: parallel processing). 

This library was opensourced so others can _not worry about it_ and just work with the metadata. :) 

## Usage

Under the covers `go-exiftool` does this: 

* `exiftool -json -binary --printConv <filename>` 
* parses the JSON into a `map[string]interface{}`
* provides helper functions to attmpt to turn `interface{}` into a typed value.


See: [GoDoc Document](https://godoc.org/github.com/mostlygeek/go-exiftool) for complete reference. 

```golang
metdadata, err := exiftool.Extract("path/to/file.JPG")

// get a string
val, err := metadata.GetString("FileName")

// get a float64, by default all numbers are float64, cause JSON
val, err := metadata.GetFloat64("GPSLongitude")

// get an int, this gets a float64 and converts it to an int
val, err := metadata.GetInt("ISO")

/*** 
 * Helpers for well known keys
 ***/

// gets the detected MIME type
mimetype := metadata.MIMEType()

// gets the "CreateDate" key as a time.Time
created, found := metadata.CreateDate()

// gets any exiftool parsing errors
exifError := metadata.Error()
```

## License 

```
MIT License

Copyright (c) 2017 Benson Wong

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```