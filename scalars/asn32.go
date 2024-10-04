package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Asn32 string

func (a *Asn32) UnmarshalGQL(v interface{}) error {

	switch v := v.(type) {
	case string:
		*a = Asn32(v)
	case int:
		*a = Asn32(strconv.Itoa(v))
	case int64:
		*a = Asn32(strconv.FormatInt(v, 10))
	case json.Number:
		srcInt, err := v.Int64()
		if err != nil {
			return fmt.Errorf("UnmarshalGQL failure to convert Asn16 to Int64: %s", v)
		}
		*a = Asn32(strconv.FormatInt(srcInt, 10))
	case nil:
		*a = Asn32("")
	default:
		*a = Asn32("")
	}

	return nil
}

func (a Asn32) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, string(a))
}

func (a Asn32) GetString() string {
	return string(a)
}
