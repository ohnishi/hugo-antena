package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/ohnishi/antena/backend/cmd"
	"github.com/pkg/errors"
)

func publishAmazonItems(use, src, dest string, date time.Time) error {
	srcJsonPath := filepath.Join(src, date.Format("20060102")+".jsonl")
	contents, err := readArticles(srcJsonPath)
	if err != nil {
		return err
	}

	feed := &feeds.Feed{
		Title:       "Amazon本アンテナ",
		Link:        &feeds.Link{Href: "http://antena.dev"},
		Description: "Amazon本RSSフィード",
		Author:      &feeds.Author{Name: "ohnishi", Email: "antena.dev@icloud.com"},
		Created:     time.Now(),
	}

	for _, content := range contents {
		t, err := time.Parse(time.RFC3339, content.Date)
		if err != nil {
			return err
		}
		item := feeds.Item{
			Title:       content.Title,
			Link:        &feeds.Link{Href: content.LinkURL},
			Description: content.Publisher,
			Content:     toContent(content),
			Created:     t,
		}
		feed.Items = append(feed.Items, &item)
	}

	f, err := cmd.CreateFile(filepath.Join(dest, "book.xml"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err = feed.WriteRss(f); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}

	return nil
}

func toContent(content cmd.AmaoznContent) string {
	var s string
	if len(content.Image) > 0 {
		s = fmt.Sprintf("<p><img src=\"%s\" /></p><p>21日放送の『バイキングMORE』（フジテレビ系）では、番組</p>", content.Image)
		if len(content.Descriptions) > 0 {
			for _, description := range content.Descriptions {
				s = s + fmt.Sprintf("<p>%s</p>", description)
			}
		}
	}
	return s
}

func readArticles(path string) ([]cmd.AmaoznContent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	var articles []cmd.AmaoznContent
	d := json.NewDecoder(f)
	for d.More() {
		var article cmd.AmaoznContent
		if err := d.Decode(&article); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal: %v", article)
		}
		articles = append(articles, article)
	}
	return articles, nil
}
