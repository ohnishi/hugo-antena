package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ohnishi/antena/backend/cmd"
	"github.com/pkg/errors"
)

func transformAmazonItems(use, src, dest string, date time.Time) (err error) {
	dateStr := date.Format("20060102")

	srcDir := filepath.Join(src, dateStr)
	amaoznLinks, err := readAmaoznLink(filepath.Join(srcDir, "link_list.jsonl"))
	if err != nil {
		return err
	}

	var acs []cmd.AmaoznContent
	for _, amaoznLink := range amaoznLinks {
		var ac cmd.AmaoznContent
		path := filepath.Join(srcDir, amaoznLink.Text+".html")

		fileInfos, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		stringReader := strings.NewReader(string(fileInfos))
		doc, err := goquery.NewDocumentFromReader(stringReader)
		if err != nil {
			panic(err)
		}

		doc.Find("img").Each(func(_ int, s *goquery.Selection) {
			_, exists := s.Attr("onload")
			if exists {
				// fmt.Println("amaoznLink.Text", amaoznLink.Text)
				imageStr, exists := s.Attr("data-a-dynamic-image")
				if exists {
					first := strings.Index(imageStr, "https://images-na.ssl-images-amazon.com/images")
					var last int
					if first >= 0 {
						last = strings.Index(imageStr[first:], "\"")
						// fmt.Println("画像", imageStr, first, last)
						ac.Image = imageStr[first : first+last]
					}
				}
			}
		})
		doc.Find("noscript").Each(func(_ int, s *goquery.Selection) {
			if !strings.Contains(s.Text(), "<em></em>") {
				return
			}
			detailStr := strings.TrimSpace(s.Text())
			first := len("<div>")
			var last int
			if first >= 0 {
				last = strings.LastIndex(detailStr, "</div>")
			}

			detailStr = strings.ReplaceAll(detailStr[first:last], "<b>", "")
			detailStr = strings.ReplaceAll(detailStr, "</b>", "")
			detailStr = strings.ReplaceAll(detailStr, "<BR>", "<br>")
			detailStr = strings.ReplaceAll(detailStr, "<h4>", "")
			detailStr = strings.ReplaceAll(detailStr, "</h4>", "")

			detailStrs := strings.Split(detailStr, "<br>")

			for _, ds := range detailStrs {
				// fmt.Println("詳細", ds)
				ds = strings.TrimSpace(ds)
				ds = strings.ReplaceAll(ds, "&emsp;", "　")
				ds = strings.ReplaceAll(ds, "&nbsp;", " ")
				ac.Descriptions = append(ac.Descriptions, ds)
			}
		})
		doc.Find("span").Each(func(_ int, s *goquery.Selection) {
			class, exists := s.Attr("class")
			if exists && class == "a-list-item" && strings.Contains(s.Text(), "出版社") {
				publisher := s.Text()
				first := strings.Index(publisher, ":")
				last := strings.LastIndex(publisher, " ")
				// fmt.Println("出版社", publisher[first+3:last])
				ac.Publisher = publisher[first+3 : last]
			}
		})

		doc.Find("span").Each(func(_ int, s *goquery.Selection) {
			class, exists := s.Attr("class")
			if exists && class == "a-list-item" {
				s.Children().Find("a").Each(func(_ int, cs *goquery.Selection) {
					href, exists := cs.Attr("href")
					if exists && strings.Contains(href, "/gp/bestsellers/books") {
						// fmt.Println("カテゴリ", strings.TrimSpace(strings.ReplaceAll(cs.Text(), " (本)", "")))
						ac.Categorys = append(ac.Categorys, strings.TrimSpace(strings.ReplaceAll(cs.Text(), " (本)", "")))
					}
				})
			}
		})
		// fmt.Println("タイトル", amaoznLink.Text)
		ac.Title = strings.ReplaceAll(amaoznLink.Text, "\"", "")

		first := strings.Index(amaoznLink.Href, "/dp/")
		first += len("/dp/")
		last := strings.Index(amaoznLink.Href[first:], "/")
		linkURL := fmt.Sprintf("https://www.amazon.co.jp/gp/product/%s", amaoznLink.Href[first:first+last])
		// fmt.Println("URL", linkURL)
		ac.LinkURL = linkURL
		ac.Date = date.Format("2006-01-02") + "T00:00:01+09:00"

		if ac.Image == "" || ac.LinkURL == "" {
			fmt.Println("skip transform title : " + ac.Title)
			continue
		}
		acs = append(acs, ac)
	}

	return writeContentAmazon(dest, dateStr+".jsonl", acs)
}

func readAmaoznLink(path string) ([]cmd.AmaoznLink, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	var links []cmd.AmaoznLink
	d := json.NewDecoder(f)
	for d.More() {
		var al cmd.AmaoznLink
		if err := d.Decode(&al); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal: %v", al)
		}
		links = append(links, al)
	}
	return links, nil
}

func writeContentAmazon(dest, fileName string, acs []cmd.AmaoznContent) error {
	f, err := cmd.CreateFile(filepath.Join(dest, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	for _, ac := range acs {
		err = cmd.AppendOutFile(f, ac)
		if err != nil {
			return err
		}
	}

	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}
