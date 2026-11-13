package scalars

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestApplicationRiskMarshalJSONNumeric(t *testing.T) {
	b, err := json.Marshal(ApplicationRisk("7"))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// The policy-platform backend rejects quoted risk scores.
	if string(b) != "7" {
		t.Fatalf("expected unquoted 7, got %s", string(b))
	}
}

func TestApplicationRiskMarshalJSONNonNumeric(t *testing.T) {
	b, err := json.Marshal(ApplicationRisk("HIGH"))
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(b) != `"HIGH"` {
		t.Fatalf("expected quoted string, got %s", string(b))
	}
}

func TestApplicationRiskMarshalGQL(t *testing.T) {
	var sb strings.Builder
	ApplicationRisk("3").MarshalGQL(&sb)
	if sb.String() != "3" {
		t.Fatalf("expected unquoted 3, got %s", sb.String())
	}
	sb.Reset()
	ApplicationRisk("HIGH").MarshalGQL(&sb)
	if sb.String() != `"HIGH"` {
		t.Fatalf("expected quoted HIGH, got %s", sb.String())
	}
}

func TestApplicationRiskUnmarshalRoundTrip(t *testing.T) {
	var r ApplicationRisk
	if err := json.Unmarshal([]byte("3"), &r); err != nil {
		t.Fatalf("unmarshal number: %v", err)
	}
	if r != "3" {
		t.Fatalf("expected \"3\", got %q", r)
	}
	b, _ := json.Marshal(r)
	if string(b) != "3" {
		t.Fatalf("round trip changed representation: %s", string(b))
	}
}

func TestApplicationRiskCanonicalNumericEncoding(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"+7", "7"}, {"007", "7"}, {"-01", "-1"}, {"-0", "0"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got, err := json.Marshal(ApplicationRisk(tc.input))
			if err != nil || string(got) != tc.want {
				t.Fatalf("JSON encoding = %q, %v; want %s", got, err, tc.want)
			}
			var gql strings.Builder
			ApplicationRisk(tc.input).MarshalGQL(&gql)
			if gql.String() != tc.want {
				t.Fatalf("GraphQL encoding = %q; want %s", gql.String(), tc.want)
			}
		})
	}
}
