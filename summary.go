package ikukani

import (
	"encoding/json"
	"time"
)

type SummaryResponse struct {
	Object        string  `json:"object"`
	URL           string  `json:"url"`
	DataUpdatedAt string  `json:"data_updated_at"`
	Data          Summary `json:"data"`
}

type Summary struct {
	NextReviewsAt string `json:"next_reviews_at"`
}

// GetSummary returns the summary struct for current user
func (c *Client) GetSummary() (*Summary, error) {
	req := request{endpoint: "/summary", method: "GET", client: c}
	resp, err := req.send()

	if err != nil {
		return nil, err
	}
	var result SummaryResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

func (c *Client) nextReviewsAt() (time.Time, error) {
	summary, err := c.GetSummary()
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.Parse(time.RFC3339, summary.NextReviewsAt)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil

}

// NextReviewIn returns the time until next review, in string format
func (c *Client) NextReviewIn() (*time.Duration, error) {
	if a, _ := c.ReviewAvailable(); a {
		d := time.Until(time.Now())
		return &d, nil
	}

	t, err := c.nextReviewsAt()
	if err != nil {
		return nil, err
	}

	d := time.Until(t)
	return &d, nil
}

// NextReviewInString returns the time until next review, in string format rounded to the minute
func (c *Client) NextReviewInString() (string, error) {
	d, err := c.NextReviewIn()
	if err != nil {
		return "", err
	}
	return d.Round(time.Minute).String(), nil
}

// ReviewAvailable checks to see if there is a current review available.
func (c *Client) ReviewAvailable() (bool, error) {
	t, err := c.nextReviewsAt()
	if err != nil {
		return false, err
	}

	n := time.Now()
	return t.Before(n), nil
}
