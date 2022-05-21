package util

import (
	wrapper "github.com/pkg/errors"
	"net/url"
	"strings"
)

func GetAbsoluteURL(domain *url.URL, u string) (string, error) {
	if strings.HasPrefix(u, "#") {
		return "", wrapper.New("url start with #")
	}
	absURL, err := domain.Parse(u)
	if err != nil {
		return "", err
	}
	absURL.Fragment = ""
	return absURL.String(), nil
}

func MustGetAbsoluteURL(domain *url.URL, u string) string {
	absoluteURL, err := GetAbsoluteURL(domain, u)
	if err != nil {
		panic(err)
	}
	return absoluteURL
}
