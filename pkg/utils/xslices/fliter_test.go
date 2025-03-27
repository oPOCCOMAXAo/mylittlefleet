package xslices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveZeroRef(t *testing.T) {
	testCases := []struct {
		input    []int
		expected []int
	}{
		{
			input:    nil,
			expected: nil,
		},
		{
			input:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			input:    []int{0, 1, 0, 2, 0, 3, 0},
			expected: []int{1, 2, 3},
		},
		{
			input:    []int{0, 0, 0, 0, 0, 0, 0},
			expected: []int{},
		},
	}
	for _, tC := range testCases {
		t.Run("", func(t *testing.T) {
			RemoveZeroRef(&tC.input)
			require.Equal(t, tC.expected, tC.input)
		})
	}
}
