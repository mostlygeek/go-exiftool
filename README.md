[![](https://godoc.org/github.com/mostlygeek/parhash?status.svg)](https://godoc.org/github.com/mostlygeek/go-exiftool)

# WARNING

I change this library a lot.  Don't consider the API stable until a version 1.0 is released

# About

go-exiftool makes it easy to extract metadata with [exiftool](https://sno.phy.queensu.ca/~phil/exiftool/) and work with it in Go.  There are currently no comparable native Go libraries with the breadth and depth of exiftool. In exchange for functionality there is a bit of a performance and a deployment penalty.

Fortunately, these are minimal. exiftool only requires perl5, which is available by default on almost every platform. The performance overhead of using an external program can be mitigated in many ways (ie: parallel processing). Using `Stayopen` and `Pool` helper libraries makes metadata extraction pretty fast.  On a 13" 2017 Macbook Pro I was able to extract metadata for 600 images in about 4 seconds.

This library was opensourced so others can _not worry about it_ and just work with the metadata. :)

## Notice

This library is still pretty young. Please use and report any bugs and issues.

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
