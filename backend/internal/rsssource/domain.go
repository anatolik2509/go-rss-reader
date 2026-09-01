package source

import "net/url"

type RssSource struct {
	Id   uint64
	Name string
	Url  url.URL
}
