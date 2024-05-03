package brave

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (b *brave) VideoSearch(ctx context.Context, term string, options ...SearchOption) (*VideoSearchResult, error) {
	u := *b.baseURL
	u.Path = u.Path + videoSearchPath

	var opts searchOptions
	applyOpts(&opts, options, nil)

	var params webSearchParams
	params.fromSearchOptions(term, opts)

	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	opts.applyRequestHeaders(b.subscriptionToken, req)

	return handleRequest[VideoSearchResult](b.client, req)
}

type VideoSearchResult struct {
	ResultContainer[VideoResult]
	Query *Query `json:"query"`
}
