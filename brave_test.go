package brave_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"dev.freespoke.com/brave-search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeb(t *testing.T) {
	svr := getTestServer("testdata/web_0.json", 200)

	client, err := brave.New("fake", brave.WithHTTPClient(svr.Client()), brave.WithBaseURL(svr.URL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.WebSearch(context.Background(), "speaker of the house")
	assert.Nil(t, err)
	if res == nil {
		t.Fatal("res was nil")
	}

	if res.Query == nil {
		t.Fatal("res.Query was nil")
	}

	assert.Equal(t, "facebook", res.Query.Original)
	assert.Equal(t, "modified", res.Query.Altered)
	assert.True(t, res.Query.ShowStrictWarning)
	assert.True(t, res.Query.Safesearch)
	assert.True(t, res.Query.IsNavigational)
	assert.True(t, res.Query.IsGeolocal)
	assert.Equal(t, "drop", res.Query.LocalDecision)
	assert.Equal(t, 1, res.Query.LocalLocationsIdx)
	assert.True(t, res.Query.IsTrending)
	assert.True(t, res.Query.IsNewsBreaking)
	assert.True(t, res.Query.AskForLocation)
	assert.Equal(t, "foo", res.Query.Language.Main)
	assert.True(t, res.Query.SpellcheckOff)
	assert.Equal(t, "us", res.Query.Country)
	assert.True(t, res.Query.BadResults)
	assert.True(t, res.Query.ShouldFallback)
	assert.Equal(t, "180.00", res.Query.Lat)
	assert.Equal(t, "-1.1", res.Query.Long)
	assert.Equal(t, "48999", res.Query.PostalCode)
	assert.Equal(t, "Detroit", res.Query.City)
	assert.Equal(t, "MI", res.Query.State)
	assert.Equal(t, "us", res.Query.HeaderCountry)
	assert.True(t, res.Query.MoreResultsAvailable)
	assert.Equal(t, "Detroit, Michigan, US", res.Query.CustomLocationLabel)

	if res.Videos == nil {
		t.Fatal("res.Videos was nil")
	}

	assert.Len(t, res.Videos.Results, 3)
	assert.True(t, res.Videos.MutatedByGoggles)
}

func TestImage(t *testing.T) {
	svr := getTestServer("testdata/images.json", 200)

	client, err := brave.New("fake", brave.WithHTTPClient(svr.Client()), brave.WithBaseURL(svr.URL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.ImageSearch(context.Background(), "speaker of the house")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestVideo(t *testing.T) {
	svr := getTestServer("testdata/videos.json", 200)

	client, err := brave.New("fake", brave.WithHTTPClient(svr.Client()), brave.WithBaseURL(svr.URL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.VideoSearch(context.Background(), "speaker of the house")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestDuration(t *testing.T) {
	var getDuration = func(in string) brave.Duration {
		d, err := time.ParseDuration(in)
		if err != nil {
			panic(err)
		}

		return brave.Duration(d)
	}

	cases := []struct {
		input    string
		expected brave.Duration
	}{
		{
			"02:04",
			getDuration("2m4s"),
		},
		{
			"59:04",
			getDuration("59m4s"),
		},
		{
			"02:02:04",
			getDuration("2h2m4s"),
		},
	}

	for _, c := range cases {
		var d brave.Duration
		assert.Nil(t, json.Unmarshal([]byte(`"`+c.input+`"`), &d))
		assert.Equal(t, c.expected, d)
	}
}

func TestRecipe(t *testing.T) {
	svr := getTestServer("testdata/web_recipe.json", 200)

	client, err := brave.New("fake", brave.WithHTTPClient(svr.Client()), brave.WithBaseURL(svr.URL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.WebSearch(context.Background(), "speaker of the house")
	require.Nil(t, err)
	require.NotNil(t, res)

	results := res.Web.Results
	require.Len(t, results, 1)
	r := results[0]
	require.NotNil(t, r.Recipe)
	require.Equal(t, "recipe", r.Subtype)

	assert.Equal(t, "Chicken Alfredo", r.Recipe.Title)
	assert.Equal(t, "desc", r.Recipe.Description)
	assert.Equal(t, 40*time.Minute, *r.Recipe.Time.Duration())
}

func getTestServer(file string, status int) *httptest.Server {
	body, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != 0 {
			w.WriteHeader(status)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		_, _ = w.Write(body)
	}))
}
