package pagination

type ForQueryParamsOutput struct {
	Limit  int
	Offset int
}

func ForQueryParams(limit int, page int) ForQueryParamsOutput {
	return ForQueryParamsOutput{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}
