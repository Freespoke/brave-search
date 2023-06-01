package brave

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

func (b *brave) WebSearch(ctx context.Context, term string, options ...SearchOption) (*WebSearchResult, error) {
	u := *b.baseURL
	u.Path = u.Path + webSearchPath

	var opts searchOptions
	applyOpts(&opts, options, nil)

	var params webSearchParams
	params.fromSearchOptions(term, opts)

	values, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	u.RawQuery = rawQuery(opts.getResultFilter(), values)

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

	var resp WebSearchResult
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type WebSearchResult struct {
	Type        string                             `json:"type"`
	Discussions *ResultContainer[DiscussionResult] `json:"discussions"`
	FAQ         any                                `json:"faq"`
	InfoBox     *ResultContainer[GraphInfoBox]     `json:"infobox"`
	Locations   any                                `json:"locations"`
	Mixed       *Mixed                             `json:"mixed"`
	News        *ResultContainer[NewsResult]       `json:"news"`
	Query       *Query                             `json:"query"`
	Videos      *ResultContainer[VideoResult]      `json:"videos"`
	Web         *ResultContainer[SearchResult]     `json:"web"`
}

type webSearchParams struct {
	Term            string `url:"q"`
	Country         string `url:"country,omitempty"`
	SearchLang      string `url:"search_lang,omitempty"`
	UILang          string `url:"ui_lang,omitempty"`
	Count           int    `url:"count,omitempty"`
	Offset          int    `url:"offset,omitempty"`
	Safesearch      string `url:"safesearch,omitempty"`
	Freshness       string `url:"freshness,omitempty"`
	TextDecorations bool   `url:"text_decorations,omitempty"`
	GogglesID       string `url:"goggles_id,omitempty"`
	Units           string `url:"units,omitempty"`
	ExtraSnippets   bool   `url:"extra_snippets,omitempty"`
}

func (w *webSearchParams) fromSearchOptions(term string, options searchOptions) {
	w.Term = term
	w.Country = options.country
	w.SearchLang = options.lang
	w.UILang = options.uiLang
	w.Count = options.count
	w.Offset = options.offset
	w.Safesearch = options.safesearch.String()
	w.Freshness = options.getFreshness()
	w.TextDecorations = options.textDecorations
	w.GogglesID = options.gogglesID
	w.Units = options.units.String()
	w.ExtraSnippets = options.extraSnippets
}
