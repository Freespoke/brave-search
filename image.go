package brave

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (b *brave) ImageSearch(ctx context.Context, term string, options ...SearchOption) (*ImageSearchResult, error) {
	u := *b.baseURL
	u.Path = u.Path + imageSearchPath

	var opts searchOptions
	applyOpts(&opts, options, func(o searchOptions) searchOptions {
		// image search does not support moderate, default to strict.
		if o.safesearch == SafesearchModerate {
			o.safesearch = SafesearchStrict
		}

		return o
	})

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

	return handleRequest[ImageSearchResult](b.client, req)
}

type ImageSearchResult struct {
	ResultContainer[ImageResult]
	Query *Query `json:"query"`
}
