package ikukani

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
)

// Token for Auth header
var Token string

// APIVersion is the wanikani api revision to use
var APIVersion = "20170710"

// Response Struct
type Response struct {
	Object        string                 `json:"object"`
	URL           string                 `json:"url"`
	DataUpdatedAt string                 `json:"data_updated_at"`
	Data          map[string]interface{} `json:"data"`
}

// Summary struct
type Summary struct {
	NextReviewsAt string `mapstructure:"next_reviews_at"`
}

// User struct
type User struct {
	CurrentVacationStartedAt string `mapstructure:"current_vacation_started_at"`
}

type request struct {
	endpoint string
	method   string
}

func (r *request) send() (*Response, error) {
	baseURL := "https://api.wanikani.com/v2/"
	if Token == "" {
		return nil, fmt.Errorf("API token needs to be set")
	}

	client := &http.Client{}
	req, err := http.NewRequest(r.method, baseURL+r.endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", `Bearer `+Token)
	req.Header.Add("Wanikani-Revision", APIVersion)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respStruct Response
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return nil, err
	}

	return &respStruct, nil
}

// GetSummary returns the summary struct for current user
func GetSummary() (*Summary, error) {
	req := request{endpoint: "/summary", method: "GET"}
	resp, err := req.send()

	if err != nil {
		return nil, err
	}
	var result Summary
	err = mapstructure.Decode(resp.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetUser returns user struct
func GetUser() (*User, error) {
	req := request{endpoint: "/user", method: "GET"}
	resp, err := req.send()

	if err != nil {
		return nil, err
	}

	var result User
	err = mapstructure.Decode(resp.Data, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func nextReviewsAt() (time.Time, error) {
	summary, err := GetSummary()
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
func NextReviewIn() (*time.Duration, error) {
	if a, _ := ReviewAvailable(); a {
		d := time.Until(time.Now())
		return &d, nil
	}

	t, err := nextReviewsAt()
	if err != nil {
		return nil, err
	}

	d := time.Until(t)
	return &d, nil
}

// NextReviewInString returns the time until next review, in string format rounded to the minute
func NextReviewInString() (string, error) {
	d, err := NextReviewIn()
	if err != nil {
		return "", err
	}
	return d.Round(time.Minute).String(), nil
}

// ReviewAvailable checks to see if there is a current review available.
func ReviewAvailable() (bool, error) {
	t, err := nextReviewsAt()
	if err != nil {
		return false, err
	}

	n := time.Now()
	return t.Before(n), nil
}

// VacationMode checks to see if vacation mode is turned on for the current user
func VacationMode() (bool, error) {
	user, err := GetUser()
	if err != nil {
		return false, err
	}
	if user.CurrentVacationStartedAt != "" {
		return true, nil
	}
	return false, nil
}
