package helpers

func NormalizePagination(page, limit int) (offset, normalizedLimit int) {
	if page <= 0 {
		page = 1
	}

	if limit <= 0 || limit > 1000 {
		limit = 5
	}

	return (page - 1) * limit, limit
}
