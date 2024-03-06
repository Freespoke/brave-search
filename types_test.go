package brave_test

import (
	"encoding/json"
	"testing"

	"dev.freespoke.com/brave-search"
	"github.com/stretchr/testify/require"
)

func TestTimestampUnmarshal(t *testing.T) {
	cases := []string{
		`"January 12, 2024"`,
		`"2024-03-06T16:41:05"`,
		`"25 minutes ago"`,
	}

	for _, c := range cases {
		var ts brave.Timestamp
		require.Nil(t, json.Unmarshal([]byte(c), &ts))
		require.False(t, ts.Time().IsZero(), c)
	}
}
