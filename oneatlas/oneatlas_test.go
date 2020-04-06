package oneatlas

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

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

func TestUnmarshallFeature(t *testing.T) {
	want_href := "https://access.foundation.api.oneatlas.airbus.com/api/v1/items/0cb0662d-f5ce-4024-83a5-0646debe34fb/images/18315065-05f0-4a5b-80c7-52790fd6f701/buffer"
	str := fmt.Sprintf(`{
  "_links": {
    "delete": {
	  "href": "https://search.foundation.api.oneatlas.airbus.com/api/v1/items/0cb0662d-f5ce-4024-83a5-0646debe34fb",
	  "name": "Delete",
	  "type": "HTTP"
	},
	"imagesGetBuffer": {
	  "href": "%s",
	  "type": "getBuffer",
	  "resourceId": "18315065-05f0-4a5b-80c7-52790fd6f701"
	}
  }
}
`, want_href)

	var f Feature
	json.NewDecoder(bytes.NewReader([]byte(str))).Decode(&f)

	log.Printf("f = %+v\n", f)

	href := f.links.ImagesGetBuffer.Links[0].Href
	if href != want_href {
		t.Errorf("href = %s; want %s", href, want_href)
	}

}

func TestUnmarshallFeatureArray(t *testing.T) {

	var fc featureCollection
	tmp, _ := ioutil.ReadFile("search.json")
	json.NewDecoder(bytes.NewReader(tmp) /*bytes.NewReader([]byte(str))*/).Decode(&fc)

	// expected failure

	log.Printf("fc = %+v\n", fc)

	// href := f.links.ImagesGetBuffer.Href
	// if href != want_href {
	// 	t.Errorf("href = %s; want %s", href, want_href)
	// }

}
