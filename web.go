package brave

import (
	"context"
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

	u.RawQuery = values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	opts.applyRequestHeaders(b.subscriptionToken, req)

	return handleRequest[WebSearchResult](b.client, req)
}

type WebSearchResult struct {
	Type                 string                             `json:"type"`
	MoreResultsAvailable bool                               `json:"more_results_available"`
	Discussions          *ResultContainer[DiscussionResult] `json:"discussions"`
	FAQ                  *ResultContainer[QA]               `json:"faq"`
	InfoBox              *ResultContainer[GraphInfoBox]     `json:"infobox"`
	Locations            *ResultContainer[LocationResult]   `json:"locations"`
	Mixed                *Mixed                             `json:"mixed"`
	News                 *ResultContainer[NewsResult]       `json:"news"`
	Query                *Query                             `json:"query"`
	Videos               *ResultContainer[VideoResult]      `json:"videos"`
	Web                  *ResultContainer[SearchResult]     `json:"web"`
	Summarizer           *Summarizer                        `json:"summarizer"`
}

type webSearchParams struct {
	Count           int      `url:"count,omitempty"`
	Country         string   `url:"country,omitempty"`
	ExtraSnippets   bool     `url:"extra_snippets,omitempty"`
	Freshness       string   `url:"freshness,omitempty"`
	GogglesID       string   `url:"goggles_id,omitempty"`
	Offset          int      `url:"offset,omitempty"`
	ResultFilter    []string `url:"result_filter,omitempty,comma"`
	Safesearch      string   `url:"safesearch,omitempty"`
	SearchLang      string   `url:"search_lang,omitempty"`
	Spellcheck      *bool    `url:"spellcheck,omitempty"`
	Term            string   `url:"q"`
	TextDecorations bool     `url:"text_decorations,omitempty"`
	UILang          string   `url:"ui_lang,omitempty"`
	Units           string   `url:"units,omitempty"`
	Summary         bool     `url:"summary,omitempty"`
}

func (w *webSearchParams) fromSearchOptions(term string, options searchOptions) {
	w.Count = options.count
	w.Country = options.country
	w.ExtraSnippets = options.extraSnippets
	w.Freshness = options.getFreshness()
	w.GogglesID = options.gogglesID
	w.Offset = options.offset
	w.ResultFilter = options.getResultFilter()
	w.Safesearch = options.safesearch.String()
	w.SearchLang = options.lang
	w.Spellcheck = options.spellcheck
	w.Term = term
	w.TextDecorations = options.textDecorations
	w.UILang = options.uiLang
	w.Units = options.units.String()
	w.Summary = options.summary
}
