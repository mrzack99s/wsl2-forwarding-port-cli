package supports

import "fmt"

func Help() {
	fmt.Print("WSL2 Forwarding Port Management \n\n")
	fmt.Println("Usage:")
	fmt.Println("	wfp-cli <command> [arguments]")
	fmt.Print("\nThe commands are:\n\n")
	fmt.Println("	ls		List all of port")
	fmt.Println("	create		Create a forwarding port rule to wsl2")
	fmt.Println("	delete		Delete a forwarding port rule to wsl2")
	fmt.Println("	purge		Purge a forwarding port rule to wsl2")
	fmt.Print("	version		Display release version\n\n")
}

func CreateHelp() {
	fmt.Print("WSL2 Forwarding Port Management \n\n")
	fmt.Println("Usage:")
	fmt.Println("	wfp-cli create --proto=<TCP|UDP> --port=<window port>:<wsl2 port> ")
}

func VersionHelp() {
	fmt.Print("WSL2 Forwarding Port Management \n\n")
	fmt.Println("Usage:")
	fmt.Println("	wfp-cli version <arguments> ")
	fmt.Print("\nThe arguments are:\n\n")
	fmt.Println("	--client		Get client version")
	fmt.Println("	--all		Get client and engine version")
}

func DeleteHelp() {
	fmt.Print("WSL2 Forwarding Port Management \n\n")
	fmt.Println("Usage:")
	fmt.Println("	wfp-cli delete <rule id> ")
}
