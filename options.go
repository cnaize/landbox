package landbox

import (
	"strconv"
	"strings"
)

type Options struct {
	TCPListen   Ports
	TCPConnect  Ports
	DenySockets bool
	DenySignals bool
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

type Ports []uint16

func (p Ports) String() string {
	ports := make([]string, 0, len(p))
	for _, port := range p {
		ports = append(ports, strconv.Itoa(int(port)))
	}

	return strings.Join(ports, ":")
}
