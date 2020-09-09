package cmds

import (
	"fmt"
	"os/exec"
	"strings"
)

func CreateRule(createArgs []string, ip string) bool {
	out, _ := exec.Command("netsh.exe", "advfirewall", "firewall", "add", "rule",
		"name="+createArgs[0]+" Port "+createArgs[1], "dir=in", "action=allow",
		"protocol="+createArgs[0], "localport="+createArgs[1]).Output()
	exec.Command("netsh.exe", "interface", "portproxy", "add", "v4tov4",
		"listenport="+createArgs[1], "listenaddress=0.0.0.0", "connectport="+createArgs[2],
		"connectaddress="+ip).Run()

	if strings.Contains(string(out), "Ok") {
		fmt.Println("Inserted success")
		return true
	}

	return false
}
