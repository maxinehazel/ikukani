package ikukani

import (
	"encoding/json"
	"fmt"
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

func getSummary() (Summary, error) {
	if Token == "" {
		return Summary{}, fmt.Errorf("API token needs to be set")
	}

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

func nextReviewsAt() (time.Time, error) {
	summary, err := getSummary()
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
func NextReviewIn() (string, error) {
	if a, _ := ReviewAvailable(); a {
		return "Review available now", nil
	}

	t, err := nextReviewsAt()
	if err != nil {
		return "", err
	}

	d := time.Until(t)
	return "Next Review in " + d.Round(time.Minute).String(), nil
}

// ReviewAvailable checks to see if there is a current review available.
func ReviewAvailable() (bool, error) {
	t, err := nextReviewsAt()
	if err != nil {
		return false, err
	}

	t = t.Round(time.Minute)
	n := time.Now().Round(time.Minute)
	return t.Before(n), nil
}
