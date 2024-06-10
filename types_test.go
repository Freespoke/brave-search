package brave_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"dev.freespoke.com/brave-search"
	"github.com/stretchr/testify/assert"
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

func TestErrorResponse(t *testing.T) {
	var resp brave.ErrorResponse
	if err := json.Unmarshal(errJSON, &resp); err != nil {
		t.Fatal(err)
	}

	p := fmt.Sprintf("%+v", resp)
	assert.Equal(t,
		"error: Unable to validate request parameter(s) (ID: f49c8ffa-5ddc-4fbf-9841-6b3093c21eb2; Status: 422; Code: VALIDATION); details: (type [int_parsing]; loc [query.offset]; input [foo]; msg [Input should be a valid integer, unable to parse string as an integer])); details: ",
		p,
	)
}

var errJSON = []byte(`{"id": "f49c8ffa-5ddc-4fbf-9841-6b3093c21eb2","status": 422,"code": "VALIDATION","detail": "Unable to validate request parameter(s)","meta": {"errors": [{"type": "int_parsing","loc": ["query","offset"],"msg": "Input should be a valid integer, unable to parse string as an integer","input": "foo"}]}}`)
