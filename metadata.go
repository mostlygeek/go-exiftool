package exiftool

import (
	"bytes"
	"encoding/base64"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

var (
	ErrKey404      = errors.New("Key does not exist")
	ErrInvalidType = errors.New("Could not convert to that type")
)

const (
	TimeFormat   = "2006:01:02 15:04:05"
	TimeFormatMS = TimeFormat + ".000"
)

// Metadata provides access to the metadata returned by exiftool
type Metadata struct {
	raw []byte
}

func NewMetadata(raw []byte) *Metadata {
	return &Metadata{raw: raw}
}

// MIMEType is a convenience function for getting the MIMEType
// key from the raw exiftool data
func (m *Metadata) MIMEType() string {
	mtype, err := jsonparser.GetString(m.raw, "File", "MIMEType")
	if err != nil {
		return ""
	} else {
		return mtype
	}
}

// CreateDate returns the date the file was created if the information
// is available
func (m *Metadata) CreateDate() (time.Time, bool) {
	var datestr string
	var err error
	datestr, err = jsonparser.GetString(m.raw, "EXIF", "CreateDate")
	if err != nil {
		datestr, err = jsonparser.GetString(m.raw, "XMP", "DateCreated")
	}

	if err != nil {
		return time.Time{}, false
	}

	t, err := time.Parse(TimeFormatMS, datestr)
	if err != nil {
		t, err = time.Parse(TimeFormat, datestr)
		if err != nil {
			return t, false
		}
	}
	return t, true
}

// GPSPosition extracts latitude, longitude from metadata. Third return value
// will be false if it was not able to find both values
func (m *Metadata) GPSPosition() (float64, float64, bool) {
	latitude, err := jsonparser.GetFloat(m.raw, "Composite", "GPSLatitude")
	if err != nil {
		return 0, 0, false
	}
	longitude, err := jsonparser.GetFloat(m.raw, "Composite", "GPSLongitude")
	if err != nil {
		return 0, 0, false
	}

	return latitude, longitude, true
}

func (m *Metadata) GetBytes(keys ...string) ([]byte, error) {
	b64, err := jsonparser.GetString(m.raw, keys...)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(b64, "base64:") {
		return nil, errors.New("does not start with `base64:`")
	}

	b64bytes, err := base64.StdEncoding.DecodeString(b64[7:])
	if err != nil {
		err = errors.Wrap(err, "Could not base64 decode JSON")
	}

	return b64bytes, err
}

// Error returns the value of `Error` key if it exists, a blank string
// otherwise. Sometimes exiftool can extract some data and still
// error because something goes wrong
func (m *Metadata) Error() string {
	str, err := jsonparser.GetString(m.raw, "ExifTool", "Error")
	if err != nil {
		if err == jsonparser.KeyPathNotFoundError {
			return ""
		} else {
			return err.Error()
		}
	}

	return str
}

func (m *Metadata) MarshalJSON() ([]byte, error) { return m.raw, nil }

// parse extracts the metadata bytes out of exiftool's output
func parse(data []byte) (*Metadata, error) {
	data = bytes.Trim(data, "[] ") // exiftool returns an array
	meta := NewMetadata(data)
	if errstr := meta.Error(); errstr != "" {
		return meta, errors.New(errstr)
	}
	return meta, nil
}
