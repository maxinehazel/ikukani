package ikukani

import (
	"encoding/json"
)

// ReviewResponse Struct
type ReviewResponse struct {
	Object        string `json:"object"`
	URL           string `json:"url"`
	DataUpdatedAt string `json:"data_updated_at"`
	Data          Review   `json:"data"`
}

type Review struct {
	AssignmentId            int    `json:"assignment_id"`
	CreatedAt               string `json:"created_at"`
	EndingSRSStageName      string `json:"ending_srs_stage_name"`
	EndingSRSStage          int    `json:"ending_srs_stage"`
	IncorrectMeaningAnswers int    `json:"incorrect_meaning_answers"`
	IncorrectReadingAnswers int    `json:"incorrect_reading_answers"`
	StartingSRSStageName    string `json:"starting_srs_stage_name"`
	StartingSRSStage        int    `json:"starting_srs_stage"`
	SubjectId               int    `json:"subject_id"`
}

func (c *Client) GetReviews() (*Review, error) {
	req := request{endpoint: "/reviews", method: "GET"}
	resp, err := req.send()

	if err != nil {
		return nil, err
	}

	var result Review
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}