package src

import (
	"encoding/json"
	"io"
)

// FromJSON deserializes the object from JSON string
// given in the io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}
