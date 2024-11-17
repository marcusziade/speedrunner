package models

// APIResponse wraps the standard API response
type APIResponse struct {
	Data []interface{} `json:"data"`
}
