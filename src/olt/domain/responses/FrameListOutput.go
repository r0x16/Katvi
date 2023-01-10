package responses

type FrameListOutput struct {
	FrameID    string `json:"frame_id"`
	FrameType  string `json:"frame_type"`
	FrameState string `json:"frame_state"`
}

type FrameListOutputCollection struct {
	FrameListOutput []FrameListOutput `json:"frame_list_output"`
	TotalFrameCount int               `json:"total_frame_count"`
}
