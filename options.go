package landbox

import (
	"strconv"
	"strings"
)

type Options struct {
	TCPListen   Ports // nil: allow all, empty: deny all
	TCPConnect  Ports // nil: allow all, empty: deny all
	DenySockets bool
	DenySignals bool
	EnableDebug bool
}

type Paths []string

func (p Paths) String() string {
	return strings.Join(p, ":")
}

type Ports []uint16

func (p Ports) String() string {
	ports := make([]string, 0, len(p))
	for _, port := range p {
		ports = append(ports, strconv.Itoa(int(port)))
	}

	return strings.Join(ports, ":")
}

func (o Options) Scope() string {
	const (
		socketsKey = "a"
		signalsKey = "s"
	)

	var scope []string
	if o.DenySockets {
		scope = append(scope, socketsKey)
	}
	if o.DenySignals {
		scope = append(scope, signalsKey)
	}

	return strings.Join(scope, ":")
}
