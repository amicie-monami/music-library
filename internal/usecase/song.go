package usecase

import (
	"fmt"
	"strings"
)

func TextPagination(text string, offset int64, limit int64) ([]string, error) {
	if limit == 0 {
		limit = 1
	}

	couplets := strings.Split(text, "\n\n")

	lenVerses := len(couplets)
	if lenVerses < int(offset) {
		return nil, fmt.Errorf("out of bounds")
	}

	if int(offset+limit) >= lenVerses {
		return couplets[offset:], nil
	}

	return couplets[offset : offset+limit], nil
}
