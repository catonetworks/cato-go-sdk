package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Long string

func (l *Long) UnmarshalGQL(v interface{}) error {

	switch v := v.(type) {
	case string:
		*l = Long(v)
	case int:
		*l = Long(strconv.Itoa(v))
	case int64:
		*l = Long(strconv.FormatInt(v, 10))
	case json.Number:
		srcInt, err := v.Int64()
		if err != nil {
			return fmt.Errorf("UnmarshalGQL failure to convert Asn16 to Int64: %s", v)
		}
		*l = Long(strconv.FormatInt(srcInt, 10))
	case nil:
		*l = Long("")
	default:
		*l = Long("")
	}

	return nil
}

func (l Long) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, string(l))
}

func (l Long) GetString() string {
	return string(l)
}
