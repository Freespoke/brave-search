package brave

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
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
	Time  int64         `json:"time"`
}

func (e *errorResponse) UnmarshalJSON(in []byte) error {
	if e == nil {
		return nil
	}

	// The Alias type is required to prevent infinite recursion back into this function.
	type Alias errorResponse

	v := &struct {
		*Alias
	}{
		(*Alias)(e),
	}

	if err := json.Unmarshal(in, v); err != nil {
		return err
	}

	e.Error.Time = time.Unix(v.Time, 0)
	*e = errorResponse(*v.Alias)

	return nil
}

type ErrorResponse struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
	Meta   struct {
		Component string `json:"component"`
	} `json:"meta"`
	Time time.Time `json:"-"`
}

func (er ErrorResponse) Error() string {
	return er.Detail
}
