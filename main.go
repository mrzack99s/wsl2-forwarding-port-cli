package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/structs"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/cmds"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/configs"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/cliparses"
	"github.com/mrzack99s/wsl2-forwarding-port-cli/supports"
)

func checkFile(filename string) error {
	_, err := os.Stat(os.Getenv("HOME") + "/." + filename)
	if os.IsNotExist(err) {
		_, err := os.Create(os.Getenv("HOME") + "/." + filename)
		os.Chmod(os.Getenv("HOME")+"/."+filename, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(filename string, rulesTable structs.RulesTable) {
	file, _ := json.MarshalIndent(rulesTable, "", " ")
	_ = ioutil.WriteFile(os.Getenv("HOME")+"/."+filename, file, 0644)
}

func asSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func checkAlreadyRuleByID(id string) bool {
	for _, rule := range configs.RulesTable.Rules {
		if id == rule.Id {
			return true
		}
	}
	return false
}

func checkAlreadyRuleBySPortAndProto(port string, proto string) bool {
	for _, rule := range configs.RulesTable.Rules {
		if port == rule.SourcePort && proto == rule.Protocol {
			return true
		}
	}
	return false
}

func main() {

	if len(os.Args) < 2 {
		supports.Help()
		os.Exit(0)
	}

	filename := "forwarding_rules.json"
	if checkFile(filename) != nil {
		fmt.Println("Json load error!!!")
		os.Exit(0)
	}

	configs.ParseForwardingTable(filename)

	var protoPtr, portPtr *string

	var ip string
	ip = ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
			}
		}
	}

	switch os.Args[1] {
	case "create":
		createCommand := flag.NewFlagSet("create", flag.ExitOnError)
		protoPtr = createCommand.String("proto", "", "Protocol <TCP|UDP>")
		portPtr = createCommand.String("port", "", "Port <window port>:<wsl2 port>")
		createCommand.Parse(os.Args[2:])
		createArgs := cliparses.CreateGetArgs(createCommand, protoPtr, portPtr)

		rule := structs.RuleStruct{
			IpAddress:       ip,
			Protocol:        createArgs[0],
			SourcePort:      createArgs[1],
			DestinationPort: createArgs[2],
		}
		hash := asSha256(rule)
		substringhash := hash[:8]

		alreadyRule := checkAlreadyRuleByID(substringhash)
		if !alreadyRule && !checkAlreadyRuleBySPortAndProto(rule.SourcePort, rule.Protocol) {
			status := cmds.CreateRule(createArgs, ip)
			if status {
				rule.Id = substringhash
				configs.RulesTable.AppendRules(rule)
				writeFile(filename, configs.RulesTable)
			}
		} else {
			if alreadyRule {
				fmt.Println("The rule is already.....")
			} else {
				fmt.Println("The source port has duplicate.....")
			}
		}
	case "delete":
		if len(os.Args) < 3 {
			supports.DeleteHelp()
			os.Exit(0)
		} else {

			id := os.Args[2]

			if checkAlreadyRuleByID(id) {
				index, status := cmds.DeleteRule(id)
				if status {
					length := len(configs.RulesTable.Rules)

					configs.RulesTable.Rules[index] = configs.RulesTable.Rules[length-1] // Copy last element to index i
					configs.RulesTable.Rules[length-1] = structs.RuleStruct{}            // Erase last element )
					configs.RulesTable.Rules = configs.RulesTable.Rules[:length-1]       // Truncate slice

					writeFile(filename, configs.RulesTable)
				}
			} else {
				fmt.Println("Not have the rule in table.....")
			}
		}

	case "ls":
		cmds.Lists()
	case "version":
		fmt.Println("WSL2-Forwarding-port-cli version 1.1.2")
	default:
		supports.Help()
		os.Exit(0)
	}

}
