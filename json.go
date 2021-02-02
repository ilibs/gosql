package gosql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var emptyJSON = json.RawMessage("{}")

func JsonObject(value interface{}) (json.RawMessage, error) {
	var source []byte
	switch t := value.(type) {
	case string:
		source = []byte(t)
	case []byte:
		source = t
	case nil:
		source = emptyJSON
	default:
		return nil, errors.New("incompatible type for json.RawMessage")
	}

	if len(source) == 0 {
		source = emptyJSON
	}

	return source, nil
}

// JSONText is a json.RawMessage, which is a []byte underneath.
// Value() validates the json format in the source, and returns an error if
// the json is not valid.  Scan does no validation.  JSONText additionally
// implements `Unmarshal`, which unmarshals the json within to an interface{}
type JSONText json.RawMessage

// MarshalJSON returns the *j as the JSON encoding of j.
func (j JSONText) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return emptyJSON, nil
	}
	return j, nil
}

// UnmarshalJSON sets *j to a copy of data
func (j *JSONText) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JSONText: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

// Value returns j as a value.  This does a validating unmarshal into another
// RawMessage.  If j is invalid json, it returns an error.
func (j JSONText) Value() (driver.Value, error) {
	var m json.RawMessage
	var err = j.Unmarshal(&m)
	if err != nil {
		return []byte{}, err
	}
	return []byte(j), nil
}

// Scan stores the src in *j.  No validation is done.
func (j *JSONText) Scan(src interface{}) error {
	var source []byte
	switch t := src.(type) {
	case string:
		source = []byte(t)
	case []byte:
		if len(t) == 0 {
			source = emptyJSON
		} else {
			source = t
		}
	case nil:
		*j = JSONText(emptyJSON)
	default:
		return errors.New("incompatible type for JSONText")
	}
	*j = append((*j)[0:0], source...)
	return nil
}

// Unmarshal unmarshal's the json in j to v, as in json.Unmarshal.
func (j *JSONText) Unmarshal(v interface{}) error {
	if len(*j) == 0 {
		*j = JSONText(emptyJSON)
	}
	return json.Unmarshal([]byte(*j), v)
}

// String supports pretty printing for JSONText types.
func (j JSONText) String() string {
	return string(j)
}

func (j *JSONText) UnmarshalBinary(data []byte) error {
	return j.UnmarshalJSON(data)
}

func (j JSONText) MarshalBinary() ([]byte, error) {
	return j.MarshalJSON()
}
