package responses

type BoardListOutput struct {
	SlotID      string `json:"slot_id"`
	BoardName   string `json:"board_name"`
	BoardStatus string `json:"board_status"`
}

type BoardListOutputCollection struct {
	BoardListOutput []*BoardListOutput `json:"board_list_output"`
	TotalBoardCount int                `json:"total_board_count"`
}
