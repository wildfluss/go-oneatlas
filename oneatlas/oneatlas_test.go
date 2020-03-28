package oneatlas

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

	features, err := client.Search()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", features)

	// Output:
	// TODO
}
