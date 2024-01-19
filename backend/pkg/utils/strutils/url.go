package strutil

import "net/url"

func AddQueryParamsToRawUrl(rawUrl string, params map[string]string) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	queries := parsedUrl.Query()
	for key, val := range params {
		queries.Add(key, val)
	}
	parsedUrl.RawQuery = queries.Encode()
	return parsedUrl.String(), nil
}
