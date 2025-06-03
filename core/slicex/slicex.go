package slicex

func Paginate[T any](slice []T, page, size int) []T {
	start := (page - 1) * size
	if start >= len(slice) {
		return []T{}
	}

	end := start + size
	if end > len(slice) {
		end = len(slice)
	}

	return slice[start:end]
}

func ToMap[K comparable, T any](rows []T, keyFunc func(row T) K) map[K]T {
	res := make(map[K]T, len(rows))
	for _, row := range rows {
		key := keyFunc(row)
		res[key] = row
	}
	return res
}
