package models

type Pusher struct {
	Subject string `json:"subject,omitempty"`
	Link    string `json:"link,omitempty"`
	Type    string `json:"type,omitempty"`
	Channel string `json:"channel,omitempty"`
	UserId  string `json:"user_id,omitempty"`
}
