package exiftool

import (
	"encoding/json"
	"time"

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
	// stores the parsed JSON data
	raw map[string]interface{}
}

func NewMetadata(raw map[string]interface{}) *Metadata {
	return &Metadata{raw: raw}
}

// MIMEType is a convenience function for getting the MIMEType
// key from the raw exiftool data
func (m *Metadata) MIMEType() string {
	if mt, err := m.GetString("MIMEType"); err == nil {
		return mt
	} else {
		return ""
	}
}

// CreateDate is a convenience function for getting the CreationDate
// as a time.Time. False is return if there was a parse error or it
// doesn't exist
func (m *Metadata) CreateDate() (time.Time, bool) {
	key := "CreateDate"
	if m.KeyExists("SubSecCreateDate") {
		key = "SubSecCreateDate"
	}

	cTime, err := m.GetDate(key)
	if err != nil {
		return cTime, false
	} else {
		return cTime, true
	}
}

// Error is a conveniece function for getting the Error key
// which is set by exiftool when it can't process something
func (m *Metadata) Error() string {
	str, err := m.GetString("Error")
	if err != nil && err != ErrKey404 {
		return err.Error()
	}

	return str
}

func (m *Metadata) GetString(key string) (string, error) {
	val, ok := m.raw[key]
	if !ok {
		return "", ErrKey404
	}

	if str, ok := val.(string); !ok {
		return "", errors.Wrap(ErrInvalidType, "Could not cast to string")
	} else {
		return str, nil
	}
}

func (m *Metadata) GetFloat64(key string) (float64, error) {
	val, ok := m.raw[key]
	if !ok {
		return 0, ErrKey404
	}

	if v, ok := val.(float64); !ok {
		return 0, errors.Wrap(ErrInvalidType, "Could not cast to float64")
	} else {
		return v, nil
	}
}

// GetInt attempts to return a value as an integer
func (m *Metadata) GetInt(key string) (int, error) {
	val, err := m.GetFloat64(key) // default for numeric types
	if err != nil {
		return 0, errors.Wrap(ErrInvalidType, "Could not get float64 to convert")
	}
	return int(val), nil
}

// GetDate formats a value as a date
func (m *Metadata) GetDate(key string) (time.Time, error) {
	str, err := m.GetString(key)
	if err != nil {
		return time.Time{}, errors.Wrap(ErrInvalidType, "Could get as string to parse")
	}

	t, err := time.Parse(TimeFormatMS, str)
	if err != nil {
		t, err = time.Parse(TimeFormat, str)
		if err != nil {
			return t, errors.Wrap(err, "Could not parse time string: "+str)
		}
	}
	return t, nil
}

// KeyExists checks if a specific value exists
func (m *Metadata) KeyExists(key string) (ok bool) {
	_, ok = m.raw[key]
	return
}

func (m *Metadata) Keys() []string {
	keys := make([]string, len(m.raw))
	var i int
	for key, _ := range m.raw {
		keys[i] = key
		i++
	}
	return keys
}

func (m *Metadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.raw)
}
