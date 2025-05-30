package logicaltest

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestNumDecodings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "basic valid input",
			input:    "12",
			expected: 2,
		},
		{
			name:     "input with zero",
			input:    "102",
			expected: 1,
		},
		{
			name:     "invalid input starting with zero",
			input:    "01",
			expected: 0,
		},
		{
			name:     "longer valid input",
			input:    "226",
			expected: 3,
		},
		{
			name:     "empty string",
			input:    "0",
			expected: 0,
		},
		{
			name:     "invalid two digit number",
			input:    "27",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, numDecodings(tt.input))
		})
	}
}
