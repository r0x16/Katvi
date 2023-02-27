package responses

type CatvPort struct {
	PortId    string `json:"id"`
	LinkState string `json:"link_state"`
	TxPower   string `json:"tx_power"`
}

type EthPort struct {
	PortId    string `json:"id"`
	PortType  string `json:"port_type"`
	Speed     string `json:"speed"`
	Duplex    string `json:"duplex"`
	LinkState string `json:"link_state"`
	RingState string `json:"ring_state"`
}

type PotsPort struct {
	PortId        string `json:"id"`
	PhysicalState string `json:"physical_state"`
	AdminState    string `json:"admin_state"`
	HookState     string `json:"hook_state"`
	SessionType   string `json:"session_type"`
	ServiceState  string `json:"service_state"`
	CallState     string `json:"call_state"`
	ServiceCodec  string `json:"service_codec"`
}

type OntLocalPortsOutputCollection struct {
	CatvPorts []*CatvPort `json:"catv_ports"`
	EthPorts  []*EthPort  `json:"eth_ports"`
	PotsPorts []*PotsPort `json:"pots_ports"`
}
