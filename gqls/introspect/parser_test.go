//nolint:goconst
package introspect

import (
	"strconv"
	"testing"
)

type tci struct {
	in, out string
}

var tcs = []tci{
	{in: "", out: ""},
	{in: "aa", out: "aa"},
	{in: "12345", out: "12345"},
	{in: "123456", out: "12345\n6"},
	{in: "1234567890", out: "12345\n67890"},
	{in: "123 567890", out: "123\n56789\n0"},
	{in: "     56789", out: "\n56789"},
	{in: "a\n\nb", out: "a\n\nb"},
}

func TestWrap(t *testing.T) {
	maxLen := 5
	for i, tc := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := wrapText(tc.in, maxLen); got != tc.out {
				t.Errorf("Wrap(%q) = %q, want %q", tc.in, got, tc.out)
			}
		})
	}
}
