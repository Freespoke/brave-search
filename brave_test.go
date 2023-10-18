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
)

func TestWeb(t *testing.T) {
	svr := getTestServer("testdata/web.json", 200)

	client, err := brave.New("fake", brave.WithHTTPClient(svr.Client()), brave.WithBaseURL(svr.URL))
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.WebSearch(context.Background(), "speaker of the house")
	assert.Nil(t, err)
	assert.NotNil(t, res)
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
