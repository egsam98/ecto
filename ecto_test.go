package ecto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrubString(t *testing.T) {
	type testCase struct {
		Input  *string
		Expect *string
	}

	for _, tc := range []testCase{
		{new("test"), new("test")},
		{new(string), nil},
		{nil, nil},
	} {
		ScrubString(&tc.Input)
		assert.Equal(t, tc.Expect, tc.Input)
	}
}
