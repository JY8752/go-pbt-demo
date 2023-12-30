package shop

import (
	"fmt"
)

func Total(items []string, prices map[string]int, _ []any) (int, error) {
	var total int
	for _, item := range items {
		p, ok := prices[item]
		if !ok {
			return 0, fmt.Errorf("not include %s in prices", item)
		}
		total += p
	}
	return total, nil
}
