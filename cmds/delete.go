package cmds

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/configs"
	"github.com/mrzack99s/wsl2-forwarding-port-cli/structs"
)

func FindElement(id string) (int, structs.RuleStruct, error) {
	for i, rule := range configs.RulesTable.Rules {
		if rule.Id == id {
			return i, rule, nil
		}
	}
	return -1, structs.RuleStruct{}, errors.New("No found")
}

func DeleteRule(id string) (int, bool) {

	index, rule, err := FindElement(id)
	if err == nil {

		exec.Command("netsh.exe", "interface", "portproxy", "delete", "v4tov4",
			"listenport="+rule.SourcePort, "listenaddress=0.0.0.0").Run()

		out, _ := exec.Command("netsh.exe", "advfirewall", "firewall", "delete", "rule",
			"name="+rule.Protocol+" Port "+rule.SourcePort, "protocol="+rule.Protocol, "localport="+rule.SourcePort).Output()

		if strings.Contains(string(out), "Ok") {
			fmt.Println("Deleted success")
			return index, true
		}
	}

	return index, false
}
