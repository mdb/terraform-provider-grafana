package gapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PauseAllAlertsResponse represents the response body for a PauseAllAlerts request.
type PauseAllAlertsResponse struct {
	AlertsAffected int64  `json:"alertsAffected,omitempty"`
	State          string `json:"state,omitempty"`
	Message        string `json:"message,omitempty"`
}

// CreateUser creates a Grafana user.
func (c *Client) CreateUser(user User) (int64, error) {
	id := int64(0)
	data, err := json.Marshal(user)
	if err != nil {
		return id, err
	}

	created := struct {
		Id int64 `json:"id"`
	}{}

	err = c.request("POST", "/api/admin/users", nil, bytes.NewBuffer(data), &created)
	if err != nil {
		return id, err
	}

	return created.Id, err
}

// DeleteUser deletes a Grafana user.
func (c *Client) DeleteUser(id int64) error {
	return c.request("DELETE", fmt.Sprintf("/api/admin/users/%d", id), nil, nil, nil)
}

// PauseAllAlerts pauses all Grafana alerts.
func (c *Client) PauseAllAlerts() (PauseAllAlertsResponse, error) {
	result := PauseAllAlertsResponse{}
	data, err := json.Marshal(PauseAlertRequest{
		Paused: true,
	})
	if err != nil {
		return result, err
	}

	err = c.request("POST", "/api/admin/pause-all-alerts", nil, bytes.NewBuffer(data), &result)

	return result, err
}
