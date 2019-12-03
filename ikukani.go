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

// Object interface is generic for data types returned
// by the wanikani api
type Object interface {
	dataType() string
}

// Summary struct
type Summary struct {
	NextReviewsAt string `mapstructure:"next_reviews_at"`
}

func (s Summary) dataType() string {
	return "report"
}

// User struct
type User struct {
	CurrentVacationStartedAt string `mapstructure:"current_vacation_started_at"`
}

func (u User) dataType() string {
	return "user"
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

func getSummary() (Summary, error) {
	req := request{endpoint: "/summary", method: "GET"}
	resp, err := req.send()

	if err != nil {
		return Summary{}, err
	}
	var result Summary
	err = mapstructure.Decode(resp.Data, &result)
	if err != nil {
		return Summary{}, err
	}

	return result, nil
}

func getUser() (User, error) {
	req := request{endpoint: "/user", method: "GET"}
	resp, err := req.send()

	if err != nil {
		return User{}, err
	}

	var result User
	err = mapstructure.Decode(resp.Data, &result)
	if err != nil {
		return User{}, err
	}

	return result, nil
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

// VacationMode checks to see if vacation mode is turned on for the current user
func VacationMode() (bool, error) {
	user, err := getUser()
	if err != nil {
		return false, err
	}
	if user.CurrentVacationStartedAt != "" {
		return true, nil
	}
	return false, nil
}
