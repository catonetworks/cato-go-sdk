package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type ID string

func (x *ID) UnmarshalGQL(v interface{}) error {

	switch v := v.(type) {
	case string:
		*x = ID(v)
	case int:
		*x = ID(strconv.Itoa(v))
	case int64:
		*x = ID(strconv.FormatInt(v, 10))
	case json.Number:
		srcInt, err := v.Int64()
		if err != nil {
			return fmt.Errorf("UnmarshalGQL failure to convert Asn16 to Int64: %s", v)
		}
		*x = ID(strconv.FormatInt(srcInt, 10))
	case nil:
		*x = ID("")
	default:
		*x = ID("")
	}

	return nil
}

func (x ID) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, string(x))
}

func (x ID) GetString() string {
	return string(x)
}
