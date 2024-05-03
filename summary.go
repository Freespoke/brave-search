package brave

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (b *brave) SummarizerSearch(ctx context.Context, key string, options ...SearchOption) (*SummarizerSearchResult, error) {
	u := *b.baseURL
	u.Path = u.Path + summarizerSearchPath

	var opts searchOptions
	applyOpts(&opts, options, nil)

	var params summarizerSearchParams
	params.fromSearchOptions(key, opts)

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

	return handleRequest[SummarizerSearchResult](b.client, req)
}

type SummarizerSearchResult struct {
	Type         string              `json:"type"`
	Status       string              `json:"status"`
	Title        string              `json:"title"`
	Summary      []SummaryMessage    `json:"summary"`
	Enrichments  *SummaryEnrichments `json:"enrichments"`
	Followups    []string            `json:"followups"`
	EntitiesInfo map[string]any      `json:"entities_info"`
}

type summarizerSearchParams struct {
	Key        string `url:"key"`
	EntityInfo bool   `url:"entity_info"`
}

func (s *summarizerSearchParams) fromSearchOptions(key string, options searchOptions) {
	s.Key = key
	s.EntityInfo = options.entityInfo
}
