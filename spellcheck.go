package brave

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (b *brave) Spellcheck(ctx context.Context, term string, options ...SearchOption) (*SpellcheckResult, error) {
	u := *b.baseURL
	u.Path = u.Path + spellcheckPath

	var opts searchOptions
	applyOpts(&opts, options, nil)

	var params spellcheckParams
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

	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var resp SpellcheckResult
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type SpellcheckResult struct {
	Type    string                 `json:"type"`
	Query   *Query                 `json:"query"`
	Results []SpellcheckResultItem `json:"results"`
}

type spellcheckParams struct {
	Term    string `url:"q"`
	Country string `url:"country,omitempty"`
	Lang    string `url:"lang,omitempty"`
}

func (s *spellcheckParams) fromSearchOptions(term string, options searchOptions) {
	s.Term = term
	s.Country = options.country
	s.Lang = options.lang
}
