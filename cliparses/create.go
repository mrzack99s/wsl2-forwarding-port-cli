package cliparses

import (
	"flag"
	"os"
	"strings"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/supports"
)

func CreateGetArgs(createCommand *flag.FlagSet, protoPtr *string, portPtr *string) []string {

	var proto, srcPort, destPort string

	if createCommand.Parsed() {

		if *protoPtr != "" {
			proto = *protoPtr
			if !(proto == "TCP" || proto == "UDP") {
				supports.CreateHelp()
				os.Exit(0)
			}
		}

		if *portPtr != "" {
			portArr := strings.Split(*portPtr, ":")
			srcPort = portArr[0]
			destPort = portArr[1]
		}

		if *portPtr == "" && *protoPtr == "" {
			supports.CreateHelp()
			os.Exit(0)
		} else {
			return []string{proto, srcPort, destPort}
		}
	}
	return nil
}
