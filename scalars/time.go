package scalars

import (
	"fmt"
	"io"
	"strconv"
)

type Time string

func (t *Time) UnmarshalGQL(v interface{}) error {
	*t = Time(v.(string))

	return nil
}

func (t Time) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(string(t)))
}

func (t Time) GetString() string {

	return string(t)
}
