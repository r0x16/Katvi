package responses

type OntListOutput struct {
	Id           string                 `json:"id"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"`
	Distance     string                 `json:"distance"`
	SerialNumber string                 `json:"serial_number"` // serial number
	RxPower      string                 `json:"rx_power"`      // rx power
	TxPower      string                 `json:"tx_power"`      // tx power
	ControlFlag  string                 `json:"control_flag"`  // control flag
	RunState     string                 `json:"run_state"`     // run state
	ConfigState  string                 `json:"config_state"`  // config state
	MatchState   string                 `json:"match_state"`   // match state
	ProtectSide  string                 `json:"protect_side"`  // protect side
	LastUpTime   string                 `json:"last_uptime"`   // last up time
	LastDownTime string                 `json:"last_downtime"` // last down time
	LastDownCase string                 `json:"last_downcase"` // last down case
	Services     map[string]*OntService `json:"services"`
}

type OntService struct {
	Index    string `json:"index"`
	VlanId   string `json:"vlan_id"`
	VlanAttr string `json:"vlan_attr"`
	PortType string `json:"port_type"`
	GemPort  string `json:"gem_port"`
	FlowType string `json:"flow_type"`
	Rx       string `json:"rx"` // rx
	Tx       string `json:"tx"` // tx
	State    string `json:"state"`
}

type OntByPonList struct {
	Onts            map[string]*OntListOutput `json:"onts"` // key: ont id
	TotalOntCount   int                       `json:"total_ont_count"`
	OnlineOntCount  int                       `json:"online_ont_count"`  // online ont count
	OfflineOntCount int                       `json:"offline_ont_count"` // offline ont count
}

type OntListOutputCollection struct {
	OntLists        map[string]*OntByPonList `json:"ont_lists"`         // key: pon id
	TotalOntCount   int                      `json:"total_ont_count"`   // total ont count
	OnlineOntCount  int                      `json:"online_ont_count"`  // online ont count
	OfflineOntCount int                      `json:"offline_ont_count"` // offline ont count
}
