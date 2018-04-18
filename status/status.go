package status

type Status struct {
	StatusPageId string `json:"statuspage_id"`
	Component    string `json:"component"`
	Container    string `json:"container"`
	Details      string `json:"details"`
	StatusID     int    `json:"current_status"`
}
