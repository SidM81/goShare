package models

type ShareRequest struct {
	Name string `json:"name"`
}

type ShareResponse struct {
	ShareURL string `json:"share_url"`
}
