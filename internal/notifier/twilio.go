package notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	TwilioSID   string
	TwilioToken string
)

// Notification is used to construct messages
type Notification struct {
	To   string
	From string
	Body string
}

// Send text messaged
func (n *Notification) Send() (string, error) {
	if TwilioSID == "" || TwilioToken == "" {
		return "", fmt.Errorf("you must provide a twilio account to send notifications")
	}

	if n.To == "" || n.From == "" || n.Body == "" {
		return "", fmt.Errorf("you must provide a To number, a From number, and a message body")
	}

	var twilioAPI = "https://api.twilio.com/2010-04-01/Accounts/" + TwilioSID + "/Messages.json"

	msgData := url.Values{}
	msgData.Set("To", n.To)
	msgData.Set("From", n.From)
	msgData.Set("Body", n.Body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", twilioAPI, &msgDataReader)
	req.SetBasicAuth(TwilioSID, TwilioToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err != nil {
			return "", err
		}
		return data["sid"].(string), nil
	}
	return "", fmt.Errorf("Error from Twilio API: " + resp.Status)

}
