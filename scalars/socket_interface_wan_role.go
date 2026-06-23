package scalars

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// SocketInterfaceWanRole is a case-insensitive enum for socket interface WAN roles.
// Some backend environments return lowercase values (e.g. "wan_1") while the canonical
// schema uses uppercase (e.g. "WAN_1"). Normalizing to uppercase on unmarshal handles both.
type SocketInterfaceWanRole string

const (
	SocketInterfaceWanRoleNone SocketInterfaceWanRole = "NONE"
	SocketInterfaceWanRoleWan1 SocketInterfaceWanRole = "WAN_1"
	SocketInterfaceWanRoleWan2 SocketInterfaceWanRole = "WAN_2"
	SocketInterfaceWanRoleWan3 SocketInterfaceWanRole = "WAN_3"
	SocketInterfaceWanRoleWan4 SocketInterfaceWanRole = "WAN_4"
)

var validSocketInterfaceWanRoles = map[SocketInterfaceWanRole]struct{}{
	SocketInterfaceWanRoleNone: {},
	SocketInterfaceWanRoleWan1: {},
	SocketInterfaceWanRoleWan2: {},
	SocketInterfaceWanRoleWan3: {},
	SocketInterfaceWanRoleWan4: {},
}

func (e SocketInterfaceWanRole) IsValid() bool {
	_, ok := validSocketInterfaceWanRoles[e]
	return ok
}

func (e SocketInterfaceWanRole) String() string {
	return string(e)
}

func (e *SocketInterfaceWanRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("SocketInterfaceWanRole must be a string")
	}
	normalized := SocketInterfaceWanRole(strings.ToUpper(str))
	if !normalized.IsValid() {
		return fmt.Errorf("%q is not a valid SocketInterfaceWanRole", str)
	}
	*e = normalized
	return nil
}

func (e SocketInterfaceWanRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(string(e)))
}

func (e *SocketInterfaceWanRole) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e SocketInterfaceWanRole) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(e))), nil
}
