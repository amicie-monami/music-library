package usecase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextPagination(t *testing.T) {
	text := `
	Ooh baby, don't you know I suffer?

	Ooh baby, can you hear me moan?

	You caught me under false pretenses
	How long before you let me go?

	Ooh
	You set my soul alight

	Ooh
	You set my soul alight
	`
	couplets, err := TextPagination(text, 0, 2)
	assert.NoError(t, err)

	fmt.Println(couplets)
}
