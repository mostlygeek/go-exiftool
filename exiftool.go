// Package exiftool provides a golang interface to the venerable exiftool
// to retrieve metadata from media files.
package exiftool

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
)

// exiftool runs exiftool and returns a decoded JSON object with
// various data types and any errors. The JSON object may be nil or not
func Extract(exiftool, filename string) (*Metadata, error) {
	// blank = exiftool on the path
	if exiftool == "" {
		exiftool = "exiftool"
	}

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
