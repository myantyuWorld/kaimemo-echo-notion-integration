package model

type KaimemoResponse struct {
	ID   string `json:"id"`
	Tag  string `json:"tag"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

type CreateKaimemoRequest struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}
