package brave

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (b *brave) SuggestSearch(ctx context.Context, term string, options ...SearchOption) (*SuggestSearchResult, error) {
	u := *b.baseURL
	u.Path = u.Path + suggestSearchPath

	var opts searchOptions
	applyOpts(&opts, options, nil)

	var params suggestParams
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

	return handleRequest[SuggestSearchResult](b.client, req)
}

type SuggestSearchResult struct {
	Type    string          `json:"type"`
	Query   *Query          `json:"query"`
	Results []SuggestResult `json:"results"`
}

type suggestParams struct {
	Term    string `url:"q"`
	Country string `url:"country,omitempty"`
	Lang    string `url:"lang,omitempty"`
	Count   int    `url:"count,omitempty"`
	Rich    bool   `url:"rich,omitempty"`
}

func (s *suggestParams) fromSearchOptions(term string, options searchOptions) {
	s.Term = term
	s.Country = options.country
	s.Lang = options.lang
	s.Count = options.count
	s.Rich = options.rich
}
