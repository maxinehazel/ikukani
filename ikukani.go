package ikukani

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Token for Auth header
var Token string

// APIVersion is the wanikani api revision to use
var APIVersion = "20170710"

// Response Struct
type Response struct {
	Object        string
	URL           string
	DataUpdatedAt string
	Data          Summary
}

// Summary struct
type Summary struct {
	NextReviewsAt string `json:"next_reviews_at"`
}

// GetSummary returns the summary data for a user
func GetSummary() (Summary, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.wanikani.com/v2/summary", nil)
	if err != nil {
		return Summary{}, err
	}

	req.Header.Add("Authorization", `Bearer `+Token)
	req.Header.Add("Wanikani-Revision", APIVersion)
	resp, err := client.Do(req)
	if err != nil {
		return Summary{}, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Summary{}, err
	}

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return Summary{}, err
	}

	return r.Data, nil
}

// NextReviewsAt for current user
func NextReviewsAt() (string, error) {
	summary, err := GetSummary()
	if err != nil {
		return "", err
	}
	t, _ := time.Parse(time.RFC3339, summary.NextReviewsAt)
	u := time.Until(t)
	return u.Round(time.Minute).String(), err
}

// ReviewAvailable checks to see if there is a current review available.
func ReviewAvailable() (bool, error) {
	summary, err := GetSummary()
	if err != nil {
		return false, err
	}
	t, err := time.Parse(time.RFC3339, summary.NextReviewsAt)
	if err != nil {
		return false, err
	}
	t = t.Round(time.Minute)
	n := time.Now().Round(time.Minute)
	return t.Before(n), nil
}
