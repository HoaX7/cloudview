package utility

type FilterParams struct {
	Page   int     `json:"page" query:"page"`
	Limit  string  `json:"limit" query:"limit"`
	Search *string `json:"search" query:"search"`
}
