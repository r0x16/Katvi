package responses

type OntListOutput struct {
	Id           string
	Description  string
	Type         string
	Distance     string
	SerialNumber string
	RxPower      string
	TxPower      string
	ControlFlag  string
	RunState     string
	ConfigState  string
	MatchState   string
	ProtectSide  string
	LastUpTime   string
	LastDownTime string
	LastDownCase string
}

type OntListOutputCollection struct {
	Onts            []*OntListOutput
	TotalOntCount   int
	OnlineOntCount  int
	OfflineOntCount int
}
