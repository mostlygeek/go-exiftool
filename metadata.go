package exiftool

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strconv"
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

// CreateDate returns the `CreateDate` key, if it exists, as a time.Time
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

// Error returns the value of `Error` key if it exists, a blank string
// otherwise. Sometimes exiftool can extract some data and still
// error because something goes wrong
func (m *Metadata) Error() string {
	str, err := m.GetString("Error")
	if err != nil && err != ErrKey404 {
		return err.Error()
	}

	return str
}

// Get returns some value that was decoded from the exiftool data
// but the type is unknown
func (m *Metadata) Get(key string) (interface{}, error) {
	val, ok := m.raw[key]
	if !ok {
		return nil, ErrKey404
	}
	return val, nil
}

// GetString returns a value as a string
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

// GetFloat64 returns a value as a float64.
func (m *Metadata) GetFloat64(key string) (float64, error) {
	val, ok := m.raw[key]
	if !ok {
		return 0, ErrKey404
	}

	if v, ok := toFloat64(val); !ok {
		return 0, errors.Wrap(ErrInvalidType, "Could not cast to float64")
	} else {
		return v, nil
	}
}

// GetFloat64s returns an []float64. Values in an array that can not be
// parsed into a float64 are returned as 0. When every value can be converted
// error will be nil
func (m *Metadata) GetFloats64s(key string) ([]float64, error) {
	val, ok := m.raw[key]
	if !ok {
		return nil, ErrKey404
	}

	if k := reflect.TypeOf(val).Kind(); k != reflect.Slice && k != reflect.Array {
		return nil, errors.New("Not an array")
	}

	s := reflect.ValueOf(val)
	num := s.Len()
	floats := make([]float64, num, num)
	allOK := true
	for i := 0; i < num; i++ {
		elem := s.Index(i).Interface()
		floats[i], allOK = toFloat64(elem)
	}

	if !allOK {
		return floats, errors.New("Could not convert all values")
	} else {
		return floats, nil
	}
}

func toFloat64(something interface{}) (float64, bool) {
	switch v := something.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case string:
		v2, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return v2, true
		} else {
			return 0, false
		}
	default:
		return 0, false
	}
}

// GetInt returns a value as an int. By default json.Unmarshal turns
// all JSON numeric types into float64 types. So that conversion is done first
// and then the float64 is turned into an int
func (m *Metadata) GetInt(key string) (int, error) {
	val, err := m.GetFloat64(key) // default for numeric types
	if err != nil {
		return 0, errors.Wrap(ErrInvalidType, "Could not get float64 to convert")
	}
	return int(val), nil
}

// GetDate parses a value into a time.Time
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

// Bytes extracts base64 encoded binary data in the metadata
func (m *Metadata) GetBytes(key string) ([]byte, error) {
	s, err := m.GetString(key)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get base64 value")
	}

	if len(s) < 8 || s[:7] != "base64:" {
		return nil, errors.New("Does not appear to be base64 encoded")
	}
	return base64.StdEncoding.DecodeString(s[7:])
}

// KeyExists checks if a specific value exists
func (m *Metadata) KeyExists(key string) (ok bool) {
	_, ok = m.raw[key]
	return
}

// Keys returns the names of all the JSON keys in the exiftool output
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
