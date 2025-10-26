package scraper

import (
	"strconv"
	"strings"
	"unicode"
)

func collapseSpaces(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func extractLines(text string) []string {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	parts := strings.Split(text, "\n")

	var lines []string
	for _, part := range parts {
		if cleaned := collapseSpaces(part); cleaned != "" {
			lines = append(lines, cleaned)
		}
	}
	return lines
}

func splitCityStateZip(line string) (string, string, string) {
	fields := strings.Fields(line)
	if len(fields) >= 3 {
		city := strings.Join(fields[:len(fields)-2], " ")
		state := fields[len(fields)-2]
		zip := fields[len(fields)-1]
		return city, state, zip
	}
	return collapseSpaces(line), "", ""
}

func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func sanitizeNumeric(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")
	return s
}

func parseFloat64(s string) float64 {
	value := sanitizeNumeric(s)
	if value == "" {
		return 0
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return f
}

func parseFloat32(s string) float32 {
	return float32(parseFloat64(s))
}

func parseInt(s string) int {
	value := sanitizeNumeric(s)
	if value == "" {
		return 0
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return n
}

func parseInt32(s string) int32 {
	return int32(parseInt(s))
}

func parseInt64(s string) int64 {
	value := sanitizeNumeric(s)
	if value == "" {
		return 0
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func appendUnique(items []string, v string) []string {
	for _, item := range items {
		if item == v {
			return items
		}
	}
	return append(items, v)
}
