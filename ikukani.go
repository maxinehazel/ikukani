package ikukani

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	Token      string
	APIVersion string
	BaseUrl    string
	Conn       *http.Client
}

func NewClient(token string, apiVersion string) *Client {
	client := Client{
		Token:      token,
		APIVersion: apiVersion,
		BaseUrl:    "https://api.wanikani.com/v2/",
		Conn:       &http.Client{},
	}

	return &client
}

type request struct {
	endpoint string
	method   string
	body     string
	client   *Client
}

func (r *request) send() ([]byte, error) {
	if r.client.Token == "" {
		return nil, fmt.Errorf("API token needs to be set")
	}

	req, err := http.NewRequest(r.method, r.client.BaseUrl+r.endpoint, strings.NewReader(r.body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", `Bearer `+r.client.Token)
	req.Header.Add("Wanikani-Revision", r.client.APIVersion)
	resp, err := r.client.Conn.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if deferErr := resp.Body.Close(); deferErr != nil {
			err = deferErr
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
