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

// NextReviewsAt for current user
func NextReviewsAt() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.wanikani.com/v2/summary", nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", `Bearer `+Token)
	req.Header.Add("Wanikani-Revision", APIVersion)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}

	t, _ := time.Parse(time.RFC3339, r.Data.NextReviewsAt)

	u := time.Until(t)
	return u.Round(time.Minute).String(), err
}
