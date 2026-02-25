package urlutils

import (
	"fmt"
	"net/url"
	"strings"
)

func Build(uri, endpoint string, params map[string]any) string {
	if strings.HasSuffix(uri, "/") {
		uri = strings.TrimSuffix(uri, "/")
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = fmt.Sprintf("%s/", endpoint)
	}
	uri = fmt.Sprintf("%s%s", uri, endpoint)
	if params != nil && len(params) > 0 {
		values := url.Values{}
		for key, value := range params {
			values.Add(key, fmt.Sprintf("%v", value))
		}
		uri = fmt.Sprintf("%s?%s", uri, values.Encode())
	}
	return uri
}
