package oneatlas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

type Link struct {
	Href string `json:"href"`
}

type links struct {
	Delete          json.RawMessage/*Link*/ `json:"delete"`
	ImagesGetBuffer json.RawMessage/*Link*/ `json:"imagesGetBuffer"`
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

	Authenticate *AuthenticateService
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
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

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
		// var err oneatlasError
		// json.NewDecoder(resp.Body).Decode(&err)
		// return nil, errors.New(err.Message)

		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, errors.New(string(bytes))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	ioutil.WriteFile("oneatlas.json", body, 0644)
	// log.Printf("%+v\n", string(body))

	err = json.NewDecoder(bytes.NewReader(body) /*resp.Body*/).Decode(v)
	return resp, err
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	u, err := url.Parse("https://search.oneatlas.geoapi-airbusds.com/")
	if err != nil {
		log.Fatal(err)
	}
	// TODO SearchService
	c := &Client{httpClient: httpClient, BaseURL: u}

	u, err = url.Parse("https://authenticate.foundation.api.oneatlas.airbus.com/")
	if err != nil {
		log.Fatal(err)
	}
	c.Authenticate = &AuthenticateService{
		c: Client{
			httpClient: httpClient,
			BaseURL:    u,
		},
	}

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

type AuthenticateService struct {
	c Client
}

type getAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	// RefreshExpiresin int `json:"refresh_expires_in"`
	// TokenType string `json:"token_type"`
}

// Get an access token.
//
// OneAtlas API docs: http://www.geoapi-airbusds.com/guides/g-authentication/
func (s *AuthenticateService) GetAccessToken(ctx context.Context, APIKey string) (string, error) {
	v := url.Values{}
	v.Set("apikey", APIKey)
	v.Set("grant_type", "api_key")
	v.Set("client_id", "IDP")

	req, err := s.c.newRequest("POST", "auth/realms/IDP/protocol/openid-connect/token", strings.NewReader(v.Encode()))
	if err != nil {
		return "", err
	}

	var resp getAccessTokenResponse
	_, err = s.c.do(ctx, req, &resp)

	return resp.AccessToken, err
}
