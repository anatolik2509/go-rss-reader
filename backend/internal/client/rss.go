package client

import (
	"context"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"rss-reader-backend/internal/service/article"
	"time"
)

type RssItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubTime     string `xml:"pubDate"`
}

type RssChannel struct {
	Title   string    `xml:"title"`
	Items   []RssItem `xml:"item"`
	Link    string    `xml:"link"`
	PubTime string    `xml:"pubDate"`
}

type RssResponse struct {
	Channel RssChannel `xml:"channel"`
}

type XmlResponse struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel RssChannel
}

type Rss struct {
}

func NewRssClient() *Rss {
	return &Rss{}
}

func (r *Rss) GetArticles(ctx context.Context, url *url.URL) ([]article.Article, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	articles, err := mapResponse(ctx, respBody)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func mapResponse(ctx context.Context, respBody []byte) (articles []article.Article, err error) {
	stringRespBody := string(respBody)
	var xmlResp XmlResponse
	err = xml.Unmarshal([]byte(stringRespBody), &xmlResp)
	if err != nil {
		return nil, err
	}
	articles = make([]article.Article, 0)
	for _, item := range xmlResp.Channel.Items {
		pubTime, err := time.Parse(time.RFC1123, item.PubTime)
		if err != nil {
			slog.WarnContext(ctx, "Failed to parse pub time: %s", item.PubTime, slog.String("error", err.Error()))
			continue
		}
		articles = append(articles, article.Article{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
			PubTime:     pubTime,
		})
	}
	return articles, nil
}
