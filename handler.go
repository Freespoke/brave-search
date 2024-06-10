package brave

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleRequest[T any](client *http.Client, req *http.Request) (*T, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in request %s: %w", req.URL.String(), err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var resp errorResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			return nil, err
		}

		resp.Error.Time = resp.Time
		resp.RawQuery = req.URL.RawQuery
		return nil, resp.Error
	}

	var resp T
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type errorResponse struct {
	Error    ErrorResponse `json:"error"`
	RawQuery string        `json:"raw_query"`
	Time     *Timestamp    `json:"time"`
}
