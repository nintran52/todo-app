package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenSalt checks if GenSalt generates a string of the correct length
func TestGenSalt(t *testing.T) {
	// Test case 1: Check if the function returns a string of the requested length
	length := 10
	result := GenSalt(length)
	assert.Equal(t, length, len(result), "GenSalt(%d) should return a string of length %d", length, length)

	// Test case 2: Check if passing a negative length returns a string of default length 50
	result = GenSalt(-1)
	assert.Equal(t, 50, len(result), "GenSalt(-1) should return a string of default length 50")

	// Test case 3: Check if passing 0 returns an empty string
	result = GenSalt(0)
	assert.Equal(t, 0, len(result), "GenSalt(0) should return an empty string")
}

// TestRandSequence checks if randSequence generates a string of the correct length
func TestRandSequence(t *testing.T) {
	// Test case 1: Check if randSequence generates a string of the correct length
	length := 15
	result := randSequence(length)
	assert.Equal(t, length, len(result), "randSequence(%d) should return a string of length %d", length, length)
}
