package random_test

import (
	"github.com/audetv/urlshortener/internal/lib/random"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewRandomString is a unit test function for NewRandomString function.
func TestNewRandomString(t *testing.T) {
	// Define test cases with different sizes
	testCases := []struct {
		name string
		size int
	}{
		{name: "size 5", size: 5},
		{name: "size 7", size: 7},
		{name: "size 10", size: 10},
		{name: "size 20", size: 20},
		{name: "size 30", size: 30},
	}

	// Iterate over the test cases
	for _, tc := range testCases {
		// Run each test case in a subtest
		t.Run(tc.name, func(t *testing.T) {
			// Generate a new random string with the given size
			str1 := random.NewRandomString(tc.size)
			// Check that the string has the expected length
			assert.Len(t, str1, tc.size)

			// Generate another random string with the same size
			str2 := random.NewRandomString(tc.size)
			// Check that the new string has the expected length
			assert.Len(t, str2, tc.size)

			// Check that the two generated strings are not equal
			assert.NotEqual(t, str1, str2)
		})
	}
}
