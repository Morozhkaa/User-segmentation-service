package models

type Segment struct {
	Slug string `json:"slug" example:"AVITO_VOICE_MESSAGES"`
}

type SegmentsList struct {
	S []string `json:"segments"`
}

type SuccessResponse struct {
	SuccessMsg string `json:"success"`
}

type UpdateRequest struct {
	SegmentsToAdd    []string `json:"segments-to-add"`
	SegmentsToRemove []string `json:"segments-to-remove"`
}
