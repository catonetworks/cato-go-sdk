package scalars

import (
	"fmt"
	"io"
)

type OperationalStatus string

func (o *OperationalStatus) UnmarshalGQL(v interface{}) error {
	*o = OperationalStatus(v.(string))

	return nil
}

func (o OperationalStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, string(o))
}

func (o OperationalStatus) GetString() string {

	return string(o)
}
