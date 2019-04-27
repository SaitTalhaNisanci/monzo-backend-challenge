package internal

import "net/url"

func absoluteUrl(currentUrl, baseUrl string) string {
	parsedUrl, _ := url.Parse(currentUrl)
	if parsedUrl.IsAbs() {
		return parsedUrl.String()
	}
	base, _ := url.Parse(baseUrl)
	return base.ResolveReference(parsedUrl).String()
}
