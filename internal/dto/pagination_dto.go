package dto

type PaginationRequest struct {
	Limit int `json:"limit" query:"limit"`
	Page  int `query:"page"`
}

type PaginationResponse struct {
	Data        interface{} `json:"data"`
	TotalItems  int64       `json:"total_items"`
	TotalPages  int64       `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	Status      int         `json:"status"`
}

func (p *PaginationRequest) SetDefaults() {
	if p.Limit == 0 {
		p.Limit = 10
	}
	if p.Page == 0 {
		p.Page = 1
	}
}
