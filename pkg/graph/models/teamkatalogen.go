package models

type TeamkatalogenResult struct {
	URL         string   `json:"url"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	NaisTeams   []string `json:"naisTeams"`
	Tags        []string `json:"tags"`
}
