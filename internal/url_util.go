package internal

import "net/url"

func absoluteUrl(currentUrl, baseUrl string) string {
	parsedUrl, err := url.Parse(currentUrl)
	if err != nil {
		return ""
	}
	if parsedUrl.IsAbs() {
		return parsedUrl.String()
	}
	base, _ := url.Parse(baseUrl)
	return base.ResolveReference(parsedUrl).String()
}
