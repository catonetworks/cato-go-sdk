package scalars

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleString unmarshals JSON strings or numbers into a string.
// Cato sometimes returns numeric literals for GraphQL String fields (e.g. file attribute sizes).
type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	if fs == nil {
		return nil
	}
	if len(data) == 0 || string(data) == "null" {
		*fs = ""
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*fs = FlexibleString(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*fs = FlexibleString(n.String())
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*fs = FlexibleString(strconv.FormatFloat(f, 'f', -1, 64))
		return nil
	}
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*fs = FlexibleString(strconv.FormatBool(b))
		return nil
	}
	return fmt.Errorf("FlexibleString: unexpected JSON %s", string(data))
}

func (fs FlexibleString) String() string {
	return string(fs)
}
