package model

type BeroepenMatch struct {
	Code   string `json:"code,omitempty"`
	Beroep string `json:"beroep,omitempty"`
}

type SearchBeroepenViaPromptRequest struct {
	Input string `json:"input,omitempty"`
}

type SearchBeroepenViaPromptResponse struct {
	Matches []BeroepenMatch `json:"matches,omitempty"`
}
