package search

import (
	"strings"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/hsm-gustavo/alnews/internal/model"
)

func Filter(items []model.Item, query string) []model.Item {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return items
	}

	result := make([]model.Item, 0)

	for _, item := range items {
		title := strings.ToLower(item.Title)

		// if its a perfect match, just add it and go to next iteration
		if strings.Contains(title, query) {
			result = append(result, item)
			continue
		}

		words := strings.Fields(title)
		best := 999

		for _, w := range words {
			// for words in title, calculate levenshtein distance and get the best number of operations
			dist := gstr.Levenshtein(w, query, 1, 1, 1)
			if dist < best {
				best = dist
			}
		}

		// if the number of operations is less than or equals to 2, add it
		if best <= 2 {
			result = append(result, item)
		}
	}

	return result
}