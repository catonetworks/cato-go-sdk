package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Asn16 string

func (a *Asn16) UnmarshalGQL(v interface{}) error {

	switch v := v.(type) {
	case string:
		*a = Asn16(v)
	case int:
		*a = Asn16(strconv.Itoa(v))
	case int64:
		*a = Asn16(strconv.FormatInt(v, 10))
	case json.Number:
		srcInt, err := v.Int64()
		if err != nil {
			return fmt.Errorf("UnmarshalGQL failure to convert Asn16 to Int64: %s", v)
		}
		*a = Asn16(strconv.FormatInt(srcInt, 10))
	case nil:
		*a = Asn16("")
	default:
		*a = Asn16("")
	}

	return nil
}

func (a Asn16) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, string(a))
}

func (a Asn16) GetString() string {
	return string(a)
}
