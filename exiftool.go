// Package exiftool provides golang bindings for calling exiftool and
// working with the metadata it is able to extract from a media file
package exiftool

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/pkg/errors"
)

// Extract calls exiftool that is available in $PATH to extract and return a
// Metadata struct. This is faster for large files like movies than ExtractReader
// since exiftool is better able to skip bytes without reading all the data.
func Extract(filename string) (*Metadata, error) {
	return ExtractCustom("exiftool", filename)
}

// ExtractCustom calls a specific exiftool executable to
// extract Metadata
func ExtractCustom(exiftool, filename string) (*Metadata, error) {
	cmd := exec.Command(exiftool, "-json", "-binary", "--printConv",
		"-groupHeadings", filename)
	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// exiftool will exit and print valid output to stdout
	// if it exits with an unrecognized filetype, don't process
	// that situtation here
	if err != nil && stdout.Len() == 0 {
		return nil, errors.Errorf("%s", stderr.String())
	}

	// no exit error but also no output
	if stdout.Len() == 0 {
		return nil, errors.New("No output")
	}

	return parse(stdout.Bytes())
}

// ExtractReader extracts metadata from an io.Reader instead of a
// filename on disk somewhere
func ExtractReader(source io.Reader) (*Metadata, error) {
	return ExtractReaderCustom("exiftool", source)
}

// ExtractReaderCustom uses a specific external exiftool to do the
// extraction
func ExtractReaderCustom(exiftool string, source io.Reader) (*Metadata, error) {
	cmd := exec.Command(exiftool, "-json", "-binary", "--printConv",
		"-groupHeadings", "-")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = source

	err := cmd.Run()

	// exiftool will exit and print valid output to stdout
	// if it exits with an unrecognized filetype, don't process
	// that situtation here
	if err != nil && stdout.Len() == 0 {
		return nil, errors.Errorf("%s", stderr.String())
	}

	// no exit error but also no output
	if stdout.Len() == 0 {
		return nil, errors.New("No output")
	}

	return parse(stdout.Bytes())
}
