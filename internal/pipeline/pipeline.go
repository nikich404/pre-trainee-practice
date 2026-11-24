package pipeline

// Double returns a new slice where each element is n*2.
// If in == nil, return nil.
func Double(in []int) []int {

	if in == nil {
		return nil
	}
	out := make([]int, len(in))
	for i, v := range in {
		out[i] = v * 2
	}
	return out
}

// FilterMin returns a new slice with elements >= min, in the same order.
// If in == nil, return nil.
func FilterMin(in []int, min int) []int {
	if in == nil {
		return nil
	}
	out := make([]int, 0, len(in))
	for _, v := range in {
		if v >= min {
			out = append(out, v)
		}
	}
	return out
}

// UniqueStable returns first occurrences only, original order kept.
// If in == nil, return nil.
func UniqueStable(in []int) []int {
	if in == nil {
		return nil
	}
	seen := make(map[int]struct{})
	out := make([]int, 0, len(in))

	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
