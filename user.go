package ikukani

import "encoding/json"

// UserResponse Struct
type UserResponse struct {
	Object        string `json:"object"`
	URL           string `json:"url"`
	DataUpdatedAt string `json:"data_updated_at"`
	Data          User   `json:"data"`
}

type User struct {
	Id                       string `json:"id"`
	Username                 string `json:"username"`
	CurrentVacationStartedAt string `json:"current_vacation_started_at"`
	Level                    int    `json:"level"`
}

// GetUser returns user struct
func (c *Client) GetUser() (*User, error) {
	req := request{endpoint: "/user", method: "GET", client: c}
	resp, err := req.send()

	if err != nil {
		return nil, err
	}

	var result UserResponse
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// VacationMode checks to see if vacation mode is turned on for the current user
func (c *Client) VacationMode() (bool, error) {
	user, err := c.GetUser()
	if err != nil {
		return false, err
	}
	if user.CurrentVacationStartedAt != "" {
		return true, nil
	}
	return false, nil
}
