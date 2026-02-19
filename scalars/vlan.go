package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Vlan int64

func (v *Vlan) UnmarshalGQL(val interface{}) error {
	switch t := val.(type) {
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			return fmt.Errorf("Vlan UnmarshalGQL: failed to convert json.Number to int64: %v", err)
		}
		*v = Vlan(i)
	case int:
		*v = Vlan(t)
	case int64:
		*v = Vlan(t)
	case float64:
		*v = Vlan(int64(t))
	case string:
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return fmt.Errorf("Vlan UnmarshalGQL: failed to parse string %q as int64: %v", t, err)
		}
		*v = Vlan(i)
	default:
		return fmt.Errorf("Vlan UnmarshalGQL: unexpected type %T", val)
	}
	return nil
}

func (v Vlan) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, int64(v))
}

func (v *Vlan) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		// Try as string
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("Vlan UnmarshalJSON: failed to unmarshal: %v", err)
		}
		parsed, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("Vlan UnmarshalJSON: failed to parse string %q: %v", s, err)
		}
		*v = Vlan(parsed)
		return nil
	}
	*v = Vlan(i)
	return nil
}

func (v Vlan) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(v))
}

func (v Vlan) GetInt64() int64 {
	return int64(v)
}
