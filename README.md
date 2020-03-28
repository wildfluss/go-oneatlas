# go-oneatlas

go-oneatlas is a Go clien library for accessing the [OneAtlas API](http://www.geoapi-airbusds.com/guides/).

## Usage ##

```go
import "github.com/ysz/go-oneatlas/oneatlas"
```

Get an access token from your API key.

```
curl -X POST https://authenticate.foundation.api.oneatlas.airbus.com/auth/realms/IDP/protocol/openid-connect/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'apikey=<api_key>&grant_type=api_key&client_id=IDP'
```

Then construct a new OneAtlas client:

```go
 
ts := oauth2.StaticTokenSource(
	&oauth2.Token{AccessToken: os.Getenv("ACCESS_TOKEN")},
)
tc := oauth2.NewClient(oauth2.NoContext, ts)

client := NewClient(tc)

```

And access different APIs. For example:

```go
// Images within a bounding box over San Francisco  
filters := &SearchFilters{Bbox: "-122.537,37.595,-122.303,37.807"}
features, err := client.Search(context.Background(), filters)
```


