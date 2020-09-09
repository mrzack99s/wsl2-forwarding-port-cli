package configs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mrzack99s/wsl2-forwarding-port-cli/structs"
)

var RulesTable structs.RulesTable
var ForwardingRulesFile *os.File

func ParseForwardingTable(filename string) {
	var err error
	ForwardingRulesFile, err = os.Open(os.Getenv("HOME") + "/." + filename)
	if err != nil {
		fmt.Println(err)
	}
	defer ForwardingRulesFile.Close()

	byteValue, _ := ioutil.ReadAll(ForwardingRulesFile)

	if string(byteValue) != "" {
		json.Unmarshal(byteValue, &RulesTable)
	}
}
