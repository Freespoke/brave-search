package brave

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiVersion        = "2023-06-01"
	defaultBaseURL    = "https://api.search.brave.com/res/v1/"
	imageSearchPath   = "images/search"
	spellcheckPath    = "spellcheck/search"
	suggestSearchPath = "suggest/search"
	videoSearchPath   = "videos/search"
	webSearchPath     = "web/search"
)

// Brave is an interface for fetching results from the Brave Search API.
type Brave interface {
	// WebSearch returns web search results.
	WebSearch(ctx context.Context, term string, options ...SearchOption) (*WebSearchResult, error)

	// SuggestSearch returns suggested related search terms.
	SuggestSearch(ctx context.Context, term string, options ...SearchOption) (*SuggestSearchResult, error)

	// Spellcheck returns spelling suggestions.
	Spellcheck(ctx context.Context, term string, options ...SearchOption) (*SpellcheckResult, error)

	// ImageSearch returns image search results.
	ImageSearch(ctx context.Context, term string, options ...SearchOption) (*ImageSearchResult, error)

	// VideoSearch returns video search results.
	VideoSearch(ctx context.Context, term string, options ...SearchOption) (*VideoSearchResult, error)
}

type brave struct {
	client            *http.Client
	baseURL           *url.URL
	subscriptionToken string
}

func New(subscriptionToken string, options ...ClientOption) (Brave, error) {
	var opts clientOptions
	applyOpts(&opts, options, func(o clientOptions) clientOptions {
		if o.baseURL == "" {
			o.baseURL = defaultBaseURL
		}

		if o.client == nil {
			o.client = http.DefaultClient
		}

		return o
	})

	u, err := url.Parse(opts.baseURL)
	if err != nil {
		return nil, err
	}

	return &brave{
		client:            opts.client,
		baseURL:           u,
		subscriptionToken: subscriptionToken,
	}, nil
}

// Freshness filters search results by when they were discovered.
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
type Freshness int8

func (f Freshness) String() string {
	switch f {
	case FreshnessPastDay:
		return "pd"
	case FreshnessPastWeek:
		return "pw"
	case FreshnessPastMonth:
		return "pm"
	case FreshnessPastYear:
		return "py"
	default:
		return ""
	}
}

// ResultFilter controls the returned data from WebSearch.
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
type ResultFilter string

func (r ResultFilter) String() string {
	switch r {
	case ResultFilterDiscussions:
		return "discussions"
	case ResultFilterFAQ:
		return "faq"
	case ResultFilterInfoBox:
		return "infobox"
	case ResultFilterNews:
		return "news"
	case ResultFilterVideos:
		return "videos"
	case ResultFilterWeb:
		return "web"
	default:
		return ""
	}
}

// Safesearch controls the adult content filter. Defaults to
// `SafesearchModerate`.
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
type Safesearch int8

func (s Safesearch) String() string {
	switch s {
	case SafesearchModerate:
		return "moderate"
	case SafesearchStrict:
		return "strict"
	default:
		return "off"
	}
}

// UnitType controls the unit of measurement used in results.
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
type UnitType int8

func (u UnitType) String() string {
	switch u {
	case UnitTypeImperial:
		return "imperial"
	case UnitTypeMetric:
		return "metric"
	default:
		return ""
	}
}

const (
	FreshnessNone Freshness = iota
	FreshnessPastDay
	FreshnessPastWeek
	FreshnessPastMonth
	FreshnessPastYear
)

const (
	ResultFilterDiscussions ResultFilter = "discussions"
	ResultFilterFAQ         ResultFilter = "faq"
	ResultFilterInfoBox     ResultFilter = "infobox"
	ResultFilterNews        ResultFilter = "news"
	ResultFilterVideos      ResultFilter = "videos"
	ResultFilterWeb         ResultFilter = "web"
	ResultFilterImages      ResultFilter = "images"
)

const (
	SafesearchModerate Safesearch = iota
	SafesearchOff
	SafesearchStrict
)

const (
	UnitTypeNone UnitType = iota
	UnitTypeMetric
	UnitTypeImperial
)

// SearchOption allows for setting optional arguments in requests.
type SearchOption func(searchOptions) searchOptions

type searchOptions struct {
	country         string
	lang            string
	uiLang          string
	count           int
	offset          int
	safesearch      Safesearch
	freshness       Freshness
	customFreshness []time.Time
	textDecorations bool
	resultFilter    []ResultFilter
	gogglesID       string
	units           UnitType
	extraSnippets   bool
	rich            bool
	noCache         bool
	userAgent       string
	locLatitude     *float32
	locLongitude    *float32
	locTimezone     *time.Location
	locCity         string
	locState        string
	locStateName    string
	locCountry      string
	locPostalCode   string
}

func (s searchOptions) getFreshness() string {
	if s.freshness == FreshnessNone && len(s.customFreshness) == 0 {
		return ""
	}

	if s.freshness != FreshnessNone {
		return s.freshness.String()
	}

	start := s.customFreshness[0]
	end := s.customFreshness[1]

	return fmt.Sprintf("%sto%s",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)
}

func (s searchOptions) getResultFilter() []string {
	if len(s.resultFilter) == 0 {
		return nil
	}

	strs := make([]string, 0, len(s.resultFilter))
	for _, r := range s.resultFilter {
		strs = append(strs, r.String())
	}

	return strs
}

func (s searchOptions) applyRequestHeaders(subscriptionToken string, req *http.Request) {
	req.Header.Add("X-Subscription-Token", subscriptionToken)

	if s.noCache {
		req.Header.Add("Cache-Control", "no-cache")
	}

	if s.userAgent != "" {
		req.Header.Add("User-Agent", s.userAgent)
	}

	if s.locLatitude != nil {
		req.Header.Add("X-Loc-Lat", fmt.Sprintf("%.3f", *s.locLatitude))
	}

	if s.locLongitude != nil {
		req.Header.Add("X-Loc-Long", fmt.Sprintf("%.3f", *s.locLongitude))
	}

	if s.locTimezone != nil {
		req.Header.Add("X-Loc-Timezone", s.locTimezone.String())
	}

	if s.locCity != "" {
		req.Header.Add("X-Loc-City", s.locCity)
	}

	if s.locState != "" {
		req.Header.Add("X-Loc-State", s.locState)
	}

	if s.locStateName != "" {
		req.Header.Add("X-Loc-State-Name", s.locStateName)
	}

	if s.locCountry != "" {
		req.Header.Add("X-Loc-Country", s.locCountry)
	}

	if s.locPostalCode != "" {
		req.Header.Add("X-Loc-Postal-Code", s.locPostalCode)
	}
}

// WithCountry specifies the search query country, where the results come from.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithCountry(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.country = v
		return o
	}
}

// WithLang specifies the search language preference.
//
// Applicable to [Brave.WebSearch] (as `search_lang`), [Brave.SuggestSearch],
// [Brave.Spellcheck].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithLang(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.lang = v
		return o
	}
}

// WithUILang specifies the user interface language preferred in response.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithUILang(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.uiLang = v
		return o
	}
}

// WithCount specifies the number of search results returned in response.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithCount(v int) SearchOption {
	return func(o searchOptions) searchOptions {
		o.count = v
		return o
	}
}

// WithOffset specifies the zero based offset that indicates number of search
// result per page (count) to skip before returning the result.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithOffset(v int) SearchOption {
	return func(o searchOptions) searchOptions {
		o.offset = v
		return o
	}
}

// WithSafesearch filters search results for adult content.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithSafesearch(v Safesearch) SearchOption {
	return func(o searchOptions) searchOptions {
		o.safesearch = v
		return o
	}
}

// WithFreshness filters search results by when they were discovered.
// To set a custom timeframe, use [WithCustomFreshness].
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithFreshness(v Freshness) SearchOption {
	return func(o searchOptions) searchOptions {
		o.freshness = v
		return o
	}
}

// WithCustomFreshness filters search results by a specified timeframe in which
// the result was discovered. To use a known value, use [WithFreshness].
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithCustomFreshness(start time.Time, end time.Time) SearchOption {
	return func(o searchOptions) searchOptions {
		o.customFreshness = []time.Time{start, end}
		return o
	}
}

// WithTextDecorations controls whether display strings, such as result
// snippets, should include decoration markers, such as highlighting characters.
// The default is true.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithTextDecorations(v bool) SearchOption {
	return func(o searchOptions) searchOptions {
		o.textDecorations = v
		return o
	}
}

// WithResultFilter specifies a list of result types to include in the search
// response.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithResultFilter(v ...ResultFilter) SearchOption {
	return func(o searchOptions) searchOptions {
		o.resultFilter = v
		return o
	}
}

// WithGogglesID specifies a goggle URL to rerank search results.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithGogglesID(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.gogglesID = v
		return o
	}
}

// WithUnits specifies the system of measurement.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithUnits(v UnitType) SearchOption {
	return func(o searchOptions) searchOptions {
		o.units = v
		return o
	}
}

// WithExtraSnippets specifies whether to return extra alternate snippets for
// web search results. Defaults to `false`.
//
// Applicable to [Brave.WebSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithExtraSnippets(v bool) SearchOption {
	return func(o searchOptions) searchOptions {
		o.extraSnippets = v
		return o
	}
}

// WithRich specifies whether to enhance suggestions with rich results. Defaults
// to `false`.
//
// Applicable to [Brave.SuggestSearch].
//
// Refer to [Query Parameters] for more detail.
//
// [Query Parameters]: https://api.search.brave.com/app/documentation/query
func WithRich(v bool) SearchOption {
	return func(o searchOptions) searchOptions {
		o.rich = v
		return o
	}
}

// WithNoCache specifies whether to disable server caching of results. Defaults
// to `false`.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithNoCache(v bool) SearchOption {
	return func(o searchOptions) searchOptions {
		o.noCache = v
		return o
	}
}

// WithUserAgent sets the user agent of the client. Defaults
// to `false`.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithUserAgent(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.userAgent = v
		return o
	}
}

// WithLocLatitude sets the latitude of the client's geographical location.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocLatitude(v float32) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locLatitude = &v
		return o
	}
}

// WithLocLongitude sets the longitude of the client's geographical location.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocLongitude(v float32) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locLongitude = &v
		return o
	}
}

// WithLocTimezone sets the timezone of the client.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocTimezone(v *time.Location) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locTimezone = v
		return o
	}
}

// WithLocCity sets the generic name of the client city.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocCity(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locCity = v
		return o
	}
}

// WithLocState sets the client state or region. Provide a two- or
// three-character value, e.g. `MI`.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocState(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locState = v
		return o
	}
}

// WithLocStateName sets the name of the client state or region.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocStateName(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locStateName = v
		return o
	}
}

// WithLocCountry sets the client country. Provide a two-letter country code,
// e.g. `US`.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocCountry(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locCountry = v
		return o
	}
}

// WithLocPostalCode sets the client postal code.
//
// Applicable to [Brave.WebSearch], [Brave.SuggestSearch], [Brave.Spellcheck].
//
// Refer to [Query Headers] for more detail.
//
// [Query Headers]: https://api.search.brave.com/app/documentation/headers
func WithLocPostalCode(v string) SearchOption {
	return func(o searchOptions) searchOptions {
		o.locPostalCode = v
		return o
	}
}

// ClientOption allows configuration of the API client.
type ClientOption func(clientOptions) clientOptions

type clientOptions struct {
	baseURL string
	client  *http.Client
}

// WithBaseURL overrides the default URL of the Brave API client.
// This is especially useful for testing.
func WithBaseURL(v string) ClientOption {
	return func(o clientOptions) clientOptions {
		o.baseURL = v
		return o
	}
}

// WithHTTPClient allows overriding of the HTTP client used to make requests to
// the Brave API.
//
// If not provided, defaults to [http.DefaultClient].
func WithHTTPClient(v *http.Client) ClientOption {
	return func(o clientOptions) clientOptions {
		o.client = v
		return o
	}
}

func applyOpts[T any, F ~func(T) T](cfg *T, opts []F, setDefaults F) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		*cfg = opt(*cfg)
	}

	if setDefaults != nil {
		*cfg = setDefaults(*cfg)
	}
}
