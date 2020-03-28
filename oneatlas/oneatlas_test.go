package oneatlas

import "context"
import "fmt"
import "golang.org/x/oauth2"
import "os"

//import "testing"

/*
func TestHelloWorld(t *testing.T) {
	want := "Hello, world."
	if got := Hello(); got != want {
		t.Errorf("Hello() = %q, want %q", got, want)
	}
}
*/

func ExampleSearch() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := NewClient(tc)

	filters := &SearchFilters{Bbox: "-122.537,37.595,-122.303,37.807"}
	features, err := client.Search(context.Background(), filters)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d images within a bounding box over San Francisco\n", len(features))

	// Output:
	// TODO
}
