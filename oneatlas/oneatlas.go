package oneatlas

import "encoding/json"
import "errors"
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
	u := c.BaseURL
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, nil
	}
	req.Header.Set("Accept", "application/json")

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

func NewClient(httpClient *http.Client) *Client {
	u, err := url.Parse("https://search.oneatlas.geoapi-airbusds.com/api/v1/opensearch")
	if err != nil {
		log.Fatal(err)
	}
	c := &Client{httpClient: httpClient, BaseURL: u}
	return c
}
