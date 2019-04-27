package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsoluteUrlWhenAbsoluteGiven(t *testing.T) {
	testCases := []struct {
		currentUrl  string
		baseUrl     string
		expectedUrl string
	}{
		{"https://monzo.com", "https://google.com", "https://monzo.com"},
		{"https://monzo.com/about", "https://google.com", "https://monzo.com/about"},
		{"https://monzo.com/careers", "https://google.com", "https://monzo.com/careers"},
		{"https://abc.com/c/d/e/f", "https://google.com", "https://abc.com/c/d/e/f"},
	}

	for _, testCase := range testCases {
		actualUrl := absoluteUrl(testCase.currentUrl, testCase.baseUrl)
		assert.Equal(t, testCase.expectedUrl, actualUrl)
	}
}

func TestAbsoluteUrlWhenRelativeGiven(t *testing.T) {
	testCases := []struct {
		currentUrl  string
		baseUrl     string
		expectedUrl string
	}{
		{"../../search", "https://monzo.com/about/a/", "https://monzo.com/search"},
		{"../search", "https://monzo.com/about/a/", "https://monzo.com/about/search"},
		{"../../../search", "https://monzo.com/about/a/", "https://monzo.com/search"},
		{"/search", "https://monzo.com/about/a/", "https://monzo.com/search"},
		{"search", "https://google.com/about/a/", "https://google.com/about/a/search"},
	}

	for _, testCase := range testCases {
		actualUrl := absoluteUrl(testCase.currentUrl, testCase.baseUrl)
		assert.Equal(t, testCase.expectedUrl, actualUrl)
	}
}
