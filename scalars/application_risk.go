package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// ApplicationRisk holds a Cato risk score, which the API may return as a JSON
// number (e.g. 3) even though the GraphQL schema declares it as a custom scalar.
// It stores the value as a string for use in Terraform state.
type ApplicationRisk string

func (r *ApplicationRisk) UnmarshalJSON(data []byte) error {
	if r == nil {
		return nil
	}
	if len(data) == 0 || string(data) == "null" {
		*r = ""
		return nil
	}
	// Try string first.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*r = ApplicationRisk(s)
		return nil
	}
	// Fall back to any JSON number.
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*r = ApplicationRisk(n.String())
		return nil
	}
	return fmt.Errorf("ApplicationRisk: unexpected JSON %s", string(data))
}

// MarshalJSON emits numeric risk values as JSON numbers — the policy-platform
// backend rejects a quoted risk score (DOWNSTREAM_SERVICE_ERROR), even though
// the GraphQL schema declares a custom scalar. Non-numeric values fall back to
// a JSON string.
func (r ApplicationRisk) MarshalJSON() ([]byte, error) {
	if _, err := strconv.ParseInt(string(r), 10, 64); err == nil {
		return []byte(string(r)), nil
	}
	return json.Marshal(string(r))
}

func (r *ApplicationRisk) UnmarshalGQL(v interface{}) error {
	switch v := v.(type) {
	case string:
		*r = ApplicationRisk(v)
	case int:
		*r = ApplicationRisk(strconv.Itoa(v))
	case int64:
		*r = ApplicationRisk(strconv.FormatInt(v, 10))
	case json.Number:
		*r = ApplicationRisk(v.String())
	case nil:
		*r = ""
	default:
		return fmt.Errorf("ApplicationRisk: unexpected type %T", v)
	}
	return nil
}

// MarshalGQL mirrors MarshalJSON: numeric values are written unquoted.
func (r ApplicationRisk) MarshalGQL(w io.Writer) {
	if _, err := strconv.ParseInt(string(r), 10, 64); err == nil {
		fmt.Fprint(w, string(r))
		return
	}
	fmt.Fprintf(w, "%q", string(r))
}

func (r ApplicationRisk) String() string {
	return string(r)
}
