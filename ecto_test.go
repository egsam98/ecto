package ecto

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestScrubString(t *testing.T) {
	type testCase struct {
		Input  *string
		Expect *string
	}

	for _, tc := range []testCase{
		{lo.ToPtr("test"), lo.ToPtr("test")},
		{new(string), nil},
		{nil, nil},
	} {
		ScrubString(&tc.Input)
		assert.Equal(t, tc.Expect, tc.Input)
	}
}
