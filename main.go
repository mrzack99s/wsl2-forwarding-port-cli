package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/cliparses"
	"github.com/mrzack99s/wsl2-forwarding-port-cli/cmds"
	"github.com/mrzack99s/wsl2-forwarding-port-cli/supports"
)

func main() {

	if len(os.Args) < 2 {
		supports.Help()
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

	var protoPtr, portPtr *string

	switch os.Args[1] {
	case "create":
		createCommand := flag.NewFlagSet("create", flag.ExitOnError)
		protoPtr = createCommand.String("proto", "", "Protocol <TCP|UDP>")
		portPtr = createCommand.String("port", "", "Port <window port>:<wsl2 port>")
		createCommand.Parse(os.Args[2:])
		createArgs := cliparses.CreateGetArgs(createCommand, protoPtr, portPtr)

		message := []byte("create@" + createArgs[0] + "@" + createArgs[1] + "@" + createArgs[2])
		_, err = udpConn.Write(message)

		for {
			buffer := make([]byte, 65535)
			n, _, err := udpConn.ReadFromUDP(buffer)
			if err != nil {
				break
			} else {
				if strings.TrimSpace(string(buffer[0:n])) == "SUCCESS" {
					fmt.Println("Inserted success")
					break
				} else if strings.TrimSpace(string(buffer[0:n])) == "ALREADY" {
					fmt.Println("The rule is already.....")
					break
				} else {
					break
				}
			}
		}

		udpConn.Close()

	case "delete":
		if len(os.Args) < 3 {
			supports.DeleteHelp()
			os.Exit(0)
		} else {

			id := os.Args[2]

			message := []byte("delete@" + id)
			_, err = udpConn.Write(message)

			for {
				buffer := make([]byte, 65535)
				n, _, err := udpConn.ReadFromUDP(buffer)
				if err != nil {
					break
				} else {
					if strings.TrimSpace(string(buffer[0:n])) == "SUCCESS" {
						fmt.Println("Deleted success")
						break
					} else if strings.TrimSpace(string(buffer[0:n])) == "ALREADY" {
						fmt.Println("The rule is already.....")
						break
					} else {
						break
					}
				}
			}

			udpConn.Close()

		}

	case "purge":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure to purge? (Y,n): ")
		ans, _ := reader.ReadString('\n')
		if strings.TrimSpace(ans) == "Y" {
			message := []byte("purge@Y")
			_, err = udpConn.Write(message)

			for {
				buffer := make([]byte, 65535)
				n, _, err := udpConn.ReadFromUDP(buffer)
				if err != nil {
					break
				} else {
					if strings.TrimSpace(string(buffer[0:n])) == "SUCCESS" {
						fmt.Println("Purge success")
						break
					} else {
						break
					}
				}
			}
		}

	case "ls":
		message := []byte("get@ls")
		_, err = udpConn.Write(message)
		bufStr := ""

		for {
			buffer := make([]byte, 65535)
			n, _, err := udpConn.ReadFromUDP(buffer)
			if err != nil {
				break
			} else {
				bufStr = strings.TrimSpace(string(buffer[0:n]))
				break
			}
		}
		cmds.Lists(bufStr)

	case "version":
		cliVersion := "2.2.0"
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
					buffer := make([]byte, 65535)
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
