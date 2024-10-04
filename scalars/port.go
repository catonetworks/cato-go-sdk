package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Port string

func (p *Port) UnmarshalGQL(v interface{}) error {
	var srcInt int64

	srcInt, err := v.(json.Number).Int64()
	if err != nil {
		return fmt.Errorf("UnmarshalGQL failure to convert to Int64: %s", v)
	}

	*p = Port(strconv.FormatInt(srcInt, 10))

	return nil
}

func (p Port) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, string(p))
}

func (p Port) GetInt64() int64 {
	strInt64, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		fmt.Println("error in GetInt64: ", err)
		return 0
	}

	return strInt64
}
