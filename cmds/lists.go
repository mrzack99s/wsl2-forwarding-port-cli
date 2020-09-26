package cmds

import (
	"fmt"
	"strings"
)

func Lists(taskListStr string) {
	fmt.Println("--------------------------------------------------------------------")
	fmt.Printf("%-10s%-22s%-12s%-12s%-12s\n", "ID", "WSL2 IPADDR", "PROTOCOL", "SPORT", "DPORT")
	fmt.Println("--------------------------------------------------------------------")
	if taskListStr != "FAILLED" {
		taskList := strings.Split(taskListStr, "@@")
		if len(taskList) > 0 {
			for _, rule := range taskList {
				ruleArgs := strings.Split(rule, "@")
				fmt.Printf("%-10s%-22s%-12s%-12s%-12s\n", ruleArgs[0], ruleArgs[1], ruleArgs[2], ruleArgs[3], ruleArgs[4])
			}
		}
	}

}
