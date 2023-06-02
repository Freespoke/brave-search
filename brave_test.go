package brave_test

import (
	"context"
	"log"
	"os"
	"testing"

	"dev.freespoke.com/brave-search"
)

func TestWeb(t *testing.T) {
	key := os.Getenv("BRAVE_API_KEY")
	if key == "" {
		t.Skip("missing BRAVE_API_KEY env")
	}

	client, err := brave.New(key)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.WebSearch(context.Background(), "facebook",
		brave.WithLocCity("Clinton Township"),
		brave.WithLocState("MI"),
		brave.WithLocCountry("US"),
		brave.WithLocPostalCode("48038"),
		brave.WithLocLatitude(42.614887),
		brave.WithLocLongitude(-82.916801),
		brave.WithSafesearch(brave.SafesearchStrict),
		brave.WithResultFilter(
			brave.ResultFilterDiscussions,
			brave.ResultFilterFAQ,
			brave.ResultFilterInfoBox,
			brave.ResultFilterNews,
			brave.ResultFilterVideos,
			brave.ResultFilterWeb,
		),
	)

	if err != nil {
		t.Fatal(err)
	}

	log.Println(res)
}
