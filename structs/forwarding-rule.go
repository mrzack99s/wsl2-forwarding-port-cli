package structs

type RulesTable struct {
	Rules []RuleStruct `json:"rules"`
}

type RuleStruct struct {
	Id              string `json:"id"`
	IpAddress       string `json:"ip_address"`
	Protocol        string `json:"protocol"`
	SourcePort      string `json:"src_port"`
	DestinationPort string `json:"dest_port"`
}

func (rulesTable *RulesTable) AppendRules(rule RuleStruct) []RuleStruct {
	rulesTable.Rules = append(rulesTable.Rules, rule)
	return rulesTable.Rules
}
