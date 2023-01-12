package responses

type PonListOutput struct {
	PortId      string `json:"port_id"`
	PortType    string `json:"port_name"`
	PortStatus  string `json:"port_status"`
	MinDistance string `json:"min_distance"`
	MaxDistance string `json:"max_distance"`
	OntTotal    string `json:"ont_total"`
	OntOnline   string `json:"ont_online"`
}

type PonListOutputCollection struct {
	PonPorts      []*PonListOutput `json:"pon_ports"`
	TotalPonCount int              `json:"total_pon_count"`
}
