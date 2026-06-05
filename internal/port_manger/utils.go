package port_manger

import (
	"fmt"
	"strings"

	commonv1 "e2b/gen/common/v1"
)

type httpRouteParams struct {
	containerID string
	domain      string
	port        *commonv1.PortInfo
}

func getPortName(port *commonv1.PortInfo) string {
	name := strings.ToLower(port.Protocol.String())
	formattedName := strings.ReplaceAll(name, "_", "-")
	return strings.ToLower(fmt.Sprintf("%s-%d", formattedName, port.PortNumber))
}

func getHTTPRoute(in httpRouteParams) string {
	return fmt.Sprintf("%s-%s.%s", in.containerID, getPortName(in.port), in.domain)
}
