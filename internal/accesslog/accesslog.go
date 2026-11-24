package accesslog

import (
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	IP     string
	Method string
	Path   string
	Status int
	Bytes  int
}

type PathCount struct {
	Path  string
	Count int
}

type Summary struct {
	Parsed    int
	Skipped   int
	UniqueIPs int
	ByStatus  map[int]int
	TopPaths  []PathCount // up to N, count desc, path asc on ties
}

// ParseLine parses one line.
// Returns ok=false for blank lines, comments, and malformed lines.
func ParseLine(s string) (Line, bool) {
	if strings.TrimSpace(s) == "" {
		return Line{}, false
	}
	if strings.HasPrefix(strings.TrimLeft(s, " \t"), "#") {
		return Line{}, false
	}

	// Split by exactly one space: double spaces give empty fields.
	fields := strings.Split(s, " ")
	if len(fields) != 5 {
		return Line{}, false
	}
	ip, method, path := fields[0], fields[1], fields[2]
	if method == "" || !strings.HasPrefix(path, "/") {
		return Line{}, false
	}
	status, err := strconv.Atoi(fields[3])
	if err != nil {
		return Line{}, false
	}
	bytes, err := strconv.Atoi(fields[4])
	if err != nil {
		return Line{}, false
	}
	return Line{IP: ip, Method: method, Path: path, Status: status, Bytes: bytes}, true
}

// Summarize parses all lines (newline-separated).
// topN is the maximum length of TopPaths. If topN <= 0, TopPaths is empty.
func Summarize(text string, topN int) Summary {
	sum := Summary{
		ByStatus: make(map[int]int),
		TopPaths: []PathCount{},
	}

	lines := strings.Split(text, "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	ips := make(map[string]struct{})
	paths := make(map[string]int)

	for _, raw := range lines {
		raw = strings.TrimSuffix(raw, "\r")
		l, ok := ParseLine(raw)
		if !ok {
			sum.Skipped++
			continue
		}
		sum.Parsed++
		ips[l.IP] = struct{}{}
		sum.ByStatus[l.Status]++
		paths[l.Path]++
	}
	sum.UniqueIPs = len(ips)

	if topN <= 0 {
		return sum
	}
	all := make([]PathCount, 0, len(paths))
	for p, c := range paths {
		all = append(all, PathCount{Path: p, Count: c})
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].Count != all[j].Count {
			return all[i].Count > all[j].Count
		}
		return all[i].Path < all[j].Path
	})
	if topN < len(all) {
		all = all[:topN]
	}
	sum.TopPaths = all
	return sum
}
