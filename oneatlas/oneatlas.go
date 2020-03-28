package oneatlas

import "context"
import "encoding/json"
import "errors"
import "fmt"
import "github.com/google/go-querystring/query"
import "io"
import "log"
import "net/http"
import "net/url"
import "reflect"
import "strings"

type Link struct {
	Href string `json:"href"`
}

type links struct {
	Delete          Link `json:"delete"`
	ImagesGetBuffer Link `json:"imagesGetBuffer"`
}

type Feature struct {
	links `json:"_links"`
}

type featureCollection struct {
	Features []Feature
}

type Client struct {
	// Base URL for API requests
	BaseURL *url.URL

	httpClient *http.Client
}

type oneatlasError struct {
	Message string `json:"message"`
}

type SearchFilters struct {
	// The search could be performed within a bounding box.
	// The box is defined by "west, south, east, north" coordinates
	// of longitude, latitude, in a EPSG:4326 decimal degrees.
	//
	// Example of a bounding box over San Francisco "-122.537,37.595,-122.303,37.807"
	Bbox string `url:"bbox,omitempty"`
}

func (c *Client) Search(ctx context.Context, filters *SearchFilters) ([]Feature, error) {
	u, err := addFilters("api/v1/opensearch", filters)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var fC featureCollection
	_, err = c.do(ctx, req, &fC)
	return fC.Features, err
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	//log.Print(u.String())
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var err oneatlasError
		json.NewDecoder(resp.Body).Decode(&err)
		return nil, errors.New(err.Message)
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

func NewClient(httpClient *http.Client) *Client {
	u, err := url.Parse("https://search.oneatlas.geoapi-airbusds.com/")
	if err != nil {
		log.Fatal(err)
	}
	c := &Client{httpClient: httpClient, BaseURL: u}
	return c
}

// h/t https://github.com/google/go-github/blob/34cb1d623f03e277545da01608448d9fea80dc3b/github/github.go#L241
func addFilters(s string, filters interface{}) (string, error) {
	v := reflect.ValueOf(filters)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(filters)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
