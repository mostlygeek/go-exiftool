package exiftool

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"github.com/bradfitz/latlong"
	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

var (
	ErrKey404      = errors.New("Key does not exist")
	ErrInvalidType = errors.New("Could not convert to that type")

	// json paths to where date information can be found in images ...
	datePaths = [][]string{
		[]string{"EXIF", "CreateDate"},
		[]string{"XMP", "DateCreated"},
	}
)

const (
	TimeFormat = "2006:01:02 15:04:05"
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
// is available. If GPS coordinates exist the return value will include
// timezone information
func (m *Metadata) CreateDate() (time.Time, bool) {
	var datestr string
	var t time.Time
	var err error

	for _, jsonPath := range datePaths {
		datestr, err = jsonparser.GetString(m.raw, jsonPath...)
		if err == nil {
			break
		}
	}

	if err != nil {
		return t, false
	}

	var location *time.Location
	if lat, long, ok := m.GPSPosition(); ok {
		if zonename := latlong.LookupZoneName(lat, long); zonename != "" {
			if loc, err := time.LoadLocation(zonename); err == nil {
				location = loc
			}
		}
	}

	if location != nil {
		t, err = time.ParseInLocation(TimeFormat, datestr, location)
	} else {
		t, err = time.Parse(TimeFormat, datestr)
	}

	return t, (err == nil)
}

// GPSPosition extracts latitude, longitude from metadata. Third return value
// will be false if it was not able to find both values
func (m *Metadata) GPSPosition() (float64, float64, bool) {

	lat, _ := jsonparser.GetString(m.raw, "Composite", "GPSLatitude")
	long, _ := jsonparser.GetString(m.raw, "Composite", "GPSLongitude")

	if lat == "" || long == "" {
		return 0, 0, false
	}

	latF, err := parseGPS(lat)
	if err != nil {
		return 0, 0, false
	}
	longF, err := parseGPS(long)
	if err != nil {
		return 0, 0, false
	}

	return latF, longF, true
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

func parseGPS(coord string) (float64, error) {
	// exiftool coords look like: 51 deg 29' 57.68" N
	parts := strings.Split(coord, " ")
	var d, m, s float64
	var negDir bool
	var err error
	for i, p := range parts {
		if i == 0 {
			d, err = strconv.ParseFloat(p, 64)
			if err != nil {
				return 0, errors.Wrap(err, "Failed parsing degrees")
			}
			continue
		}

		if p == "deg" {
			continue
		}

		if p[len(p)-1] == '\'' {
			m, err = strconv.ParseFloat(p[:len(p)-1], 64)
			if err != nil {
				return 0, errors.Wrap(err, "Failed parsing minutes")
			}
			continue
		}

		if p[len(p)-1] == '"' {
			s, err = strconv.ParseFloat(p[:len(p)-1], 64)
			if err != nil {
				return 0, errors.Wrap(err, "Failed parsing seconds")
			}
			continue
		}

		if p == "S" || p == "W" {
			negDir = true
		}
	}

	v := d + m/60.0 + s/60/60
	if negDir {
		v = 0 - v
	}

	return v, nil

}
