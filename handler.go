package brave

import (
	"context"
	"encoding/json"
	"net/http"
)

func handleRequest[T any](ctx context.Context, client *http.Client, req *http.Request) (*T, error) {
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var resp errorResponse
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			return nil, err
		}

		resp.Error.Time = resp.Time
		return nil, resp.Error
	}

	var resp T
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type errorResponse struct {
	Error ErrorResponse `json:"error"`
	Time  *Timestamp    `json:"time"`
}
