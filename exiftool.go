// Package exiftool provides a golang interface to the venerable exiftool
// to retrieve metadata from media files.
package exiftool

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
)

// Extract calls exiftool that is in the path to extract and return a
// Metadata struct
func Extract(filename string) (*Metadata, error) {
	return ExtractCustom("exiftool", filename)
}

// ExtractCustom calls a specific external `exiftool` executable to
// extract Metadata
func ExtractCustom(exiftool, filename string) (*Metadata, error) {
	cmd := exec.Command(exiftool, "-json", filename)
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

	container := make([]map[string]interface{}, 1, 1)
	err = json.Unmarshal(stdout.Bytes(), &container)
	if err != nil {
		return nil, errors.Wrap(err, "JSON unmarshal failed")
	}

	if len(container) != 1 {
		return nil, errors.New("Expected one record")
	}

	meta := NewMetadata(container[0])
	if errstr := meta.Error(); errstr != "" {
		return meta, errors.New(errstr)
	}

	return meta, nil
}
