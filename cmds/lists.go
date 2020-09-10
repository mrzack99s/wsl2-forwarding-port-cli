package cmds

import (
	"fmt"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/configs"
)

func Lists() {
	fmt.Println("--------------------------------------------------------------------")
	fmt.Printf("%-10s%-22s%-12s%-12s%-12s\n", "ID", "WSL2 IPADDR", "PROTOCOL", "SPORT", "DPORT")
	fmt.Println("--------------------------------------------------------------------")
	if len(configs.RulesTable.Rules) > 0 {
		for _, rule := range configs.RulesTable.Rules {
			fmt.Printf("%-10s%-22s%-10s%-10s%-10s\n", rule.Id, rule.IpAddress, rule.Protocol, rule.SourcePort, rule.DestinationPort)
		}
	}

}
