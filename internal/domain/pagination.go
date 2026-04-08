package domain

type Pagination struct {
	Limit int `json:"limit" query:"limit"`
	Page  int `json:"page" query:"page"`
}
