package oneatlas

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
)

//import "testing"

/*
func TestHelloWorld(t *testing.T) {
	want := "Hello, world."
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
*/

type oneatlasTokenSource struct {
	APIKey string

	client *Client
}

func (ts oneatlasTokenSource) Token() (*oauth2.Token, error) {
	accessToken, _ := ts.client.Authenticate.GetAccessToken(context.Background(), ts.APIKey)

	return &oauth2.Token{
		AccessToken: accessToken,
	}, nil
}

func TokenSource(APIKey string) oauth2.TokenSource {
	if APIKey == "" {
		panic("TokenSource: APIKey = \"\"; OneAtlas API key is required")
	}

	client := NewClient(nil)

	return oneatlasTokenSource{
		APIKey: APIKey,
		client: client,
	}
}

func ExampleSearch() {
	ts := TokenSource(os.Getenv("APIKEY"))
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := NewClient(tc)

	filters := &SearchFilters{Bbox: "-122.537,37.595,-122.303,37.807"}
	features, err := client.Search(context.Background(), filters)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%d images within a bounding box over San Francisco\n", len(features))
	log.Printf("%+v\n", features)

	// Output:
	// TODO
}

func ExampleGetAccessToken() {
	client := NewClient(nil)

	accessToken, err := client.Authenticate.GetAccessToken(context.Background(), os.Getenv("APIKEY"))

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("accessToken = %s\n", accessToken)

	// Output:
	// TODO
}
