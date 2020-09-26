package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"

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

func FindElement(id string) (int, structs.RuleStruct, error) {
	for i, rule := range configs.RulesTable.Rules {
		if rule.Id == id {
			return i, rule, nil
		}
	}
	return -1, structs.RuleStruct{}, errors.New("No found")
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

	out, _ := exec.Command("bash", "-c", "ip route | grep default | awk '{print $3}'").Output()
	winIp := strings.TrimSpace(string(out))

	udpAddr, err := net.ResolveUDPAddr("udp", winIp+":40123")
	if err != nil {
		fmt.Println("Wrong Address")
		return
	}

	//Create the connection
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println(err)
		return
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

			message := []byte("create@" + substringhash + "@" + rule.Protocol + "@" + rule.SourcePort + "@" + rule.DestinationPort)
			_, err = udpConn.Write(message)

			status := false
			//Keep calling this function
			for {
				buffer := make([]byte, 2048)
				n, _, err := udpConn.ReadFromUDP(buffer)
				if err != nil {
					break
				} else {
					if strings.TrimSpace(string(buffer[0:n])) == "SUCCESS" {
						status = true
						break
					} else if strings.TrimSpace(string(buffer[0:n])) == "ALREADY" {
						fmt.Println("The rule is already.....")
						break
					} else {
						break
					}
				}
			}

			if status {
				rule.Id = substringhash
				configs.RulesTable.AppendRules(rule)
				writeFile(filename, configs.RulesTable)
				fmt.Println("Inserted success")
			}

			udpConn.Close()

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
				index, _, err := FindElement(id)
				if err == nil {

					message := []byte("delete@" + id)
					_, err = udpConn.Write(message)

					status := false
					//Keep calling this function
					for {
						buffer := make([]byte, 2048)
						n, _, err := udpConn.ReadFromUDP(buffer)
						if err != nil {
							break
						} else {
							if strings.TrimSpace(string(buffer[0:n])) == "SUCCESS" {
								status = true
								break
							} else if strings.TrimSpace(string(buffer[0:n])) == "ALREADY" {
								fmt.Println("The rule is already.....")
								break
							} else {
								break
							}
						}
					}

					if status {
						length := len(configs.RulesTable.Rules)

						configs.RulesTable.Rules[index] = configs.RulesTable.Rules[length-1] // Copy last element to index i
						configs.RulesTable.Rules[length-1] = structs.RuleStruct{}            // Erase last element )
						configs.RulesTable.Rules = configs.RulesTable.Rules[:length-1]       // Truncate slice
						writeFile(filename, configs.RulesTable)
						fmt.Println("Delete success")
					}

					udpConn.Close()
				}
			} else {
				fmt.Println("Not have the rule in table.....")
			}
		}

	case "purge":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure to purge? (Y,n): ")
		ans, _ := reader.ReadString('\n')
		if strings.TrimSpace(ans) == "Y" {
			message := []byte("purge@Y")
			_, err = udpConn.Write(message)

			status := false
			//Keep calling this function
			for {
				buffer := make([]byte, 2048)
				n, _, err := udpConn.ReadFromUDP(buffer)
				if err != nil {
					break
				} else {
					if strings.TrimSpace(string(buffer[0:n])) == "SUCCESS" {
						status = true
						break
					} else {
						break
					}
				}
			}
			if status {
				configs.RulesTable.Rules = nil
				writeFile(filename, configs.RulesTable)
				fmt.Println("Purge success")
			}

		}

	case "ls":
		cmds.Lists()
	case "version":
		cliVersion := "2.1.0"
		if len(os.Args) < 3 {
			supports.VersionHelp()
		} else {
			if os.Args[2] == "--client" {
				fmt.Println("WSL2-Forwarding-port-cli version " + cliVersion)
			} else if os.Args[2] == "--all" {
				message := []byte("get@engine_version")
				_, err = udpConn.Write(message)

				status := false
				//Keep calling this function
				engineVersion := ""
				for {
					buffer := make([]byte, 2048)
					n, _, err := udpConn.ReadFromUDP(buffer)
					if err != nil {
						break
					} else {
						status = true
						engineVersion = strings.TrimSpace(string(buffer[0:n]))
						break
					}
				}
				if status {
					fmt.Println("WSL2-Forwarding-port-cli version " + cliVersion)
					fmt.Println("WSL2-Forwarding-port-engine version " + engineVersion)
				}
			} else {
				supports.VersionHelp()
			}
		}
	default:
		supports.Help()
		os.Exit(0)
	}

}
