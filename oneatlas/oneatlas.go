package oneatlas

import "encoding/json"
import "errors"
import "io"
import "log"
import "net/http"
import "net/url"

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

func (c *Client) Search() ([]Feature, error) {
	req, err := c.newRequest("GET", "/api/v1/opensearch", nil)
	if err != nil {
		return nil, nil
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var err oneatlasError
		json.NewDecoder(resp.Body).Decode(&err)
		return nil, errors.New(err.Message)
	}

	var fC featureCollection
	err = json.NewDecoder(resp.Body).Decode(&fC)
	return fC.Features, err
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	return req, nil
}

func NewClient(httpClient *http.Client) *Client {
	u, err := url.Parse("https://search.oneatlas.geoapi-airbusds.com")
	if err != nil {
		log.Fatal(err)
	}
	c := &Client{httpClient: httpClient, BaseURL: u}
	return c
}
