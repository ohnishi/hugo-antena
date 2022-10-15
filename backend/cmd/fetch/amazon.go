package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ohnishi/antena/backend/cmd"
	"github.com/pkg/errors"
)

var amazonURLMap = map[string]string{
	"itbook":       "https://www.amazon.co.jp/s?i=stripbooks&bbn=466298&rh=n%3A465392%2Cn%3A465610%2Cn%3A466298%2Cp_n_publication_date%3A2285919051%2Cp_n_binding_browse-bin%3A2148405051%7C86137051%7C86138051&s=date-asc-rank&dc&__mk_ja_JP=%E3%82%AB%E3%82%BF%E3%82%AB%E3%83%8A&qid=1603449624&rnid=625256011&ref=sr_st_date-asc-rank",
	"itmagazine":   "https://www.amazon.co.jp/s?i=stripbooks&bbn=46423011&rh=n%3A465392%2Cn%3A465610%2Cn%3A13384021%2Cn%3A46423011%2Cp_n_publication_date%3A2285919051&s=date-asc-rank&dc&__mk_ja_JP=%E3%82%AB%E3%82%BF%E3%82%AB%E3%83%8A&qid=1603449590&rnid=82836051&ref=sr_st_date-asc-rank",
	"mensmagazine": "https://www.amazon.co.jp/s?i=stripbooks&bbn=13384021&rh=n%3A465392%2Cn%3A13384021%2Cn%3A46429011%2Cp_n_publication_date%3A2285919051&s=date-asc-rank&dc&fst=as%3Aoff&qid=1608538334&rnid=13384021&ref=sr_nr_n_11",
	"finance":      "https://www.amazon.co.jp/s?i=stripbooks&bbn=492054&rh=n%3A465392%2Cn%3A465610%2Cn%3A492054%2Cp_n_publication_date%3A2285919051%2Cp_n_binding_browse-bin%3A2148405051%7C86137051%7C86138051%7C86140051%7C86142051&s=date-asc-rank&dc&fst=as%3Aoff&qid=1605601234&rnid=625256011&ref=sr_nr_p_n_binding_browse-bin_6",
	"health":       "https://www.amazon.co.jp/s?i=stripbooks&bbn=2133603051&rh=n%3A465392%2Cn%3A466304%2Cn%3A2133603051%2Cp_n_publication_date%3A2285919051%2Cp_n_binding_browse-bin%3A2148405051%7C86137051%7C86138051%7C86139051%7C86140051%7C86142051&s=date-asc-rank&dc&fst=as%3Aoff&qid=1605601520&rnid=625256011&ref=sr_nr_p_n_binding_browse-bin_7",
	"business":     "https://www.amazon.co.jp/s?i=stripbooks&bbn=465610&rh=n%3A465392%2Cn%3A466282%2Cp_n_publication_date%3A2285919051%2Cp_n_binding_browse-bin%3A2148405051%7C86137051%7C86138051%7C86140051%7C86142051&s=date-asc-rank&dc&fst=as%3Aoff&qid=1608538209&rnid=465610&ref=sr_nr_n_6",
}

func fetchAmazon(use, dest string, maxRetry uint) (err error) {
	amazonLinks := make(map[cmd.AmaoznLink]struct{})
	for _, amazonURL := range amazonURLMap {
		retry := uint(0)
		var links []cmd.AmaoznLink
		for {
			links, err = func() ([]cmd.AmaoznLink, error) {
				res, err := http.Get(amazonURL)
				if err != nil {
					return nil, errors.Wrap(err, "failed request amazon it itbook list")
				}
				defer res.Body.Close()

				if res.StatusCode != http.StatusOK {
					return nil, errors.Errorf("status code expected 200 but was %d", res.StatusCode)
				}

				// body, error := ioutil.ReadAll(res.Body)
				// if error != nil {
				// 	log.Fatal(error)
				// }
				// fmt.Println("[body] " + string(body))
				// if len(body) > 0 {
				// 	return nil, nil
				// }

				links, err := getLinks(res.Body)
				if err != nil {
					return nil, err
				}
				return links, nil
			}()
			retry++
			if err == nil || retry > maxRetry {
				break
			}
			time.Sleep(3 * time.Second)
		}
		if err != nil {
			return err
		}
		for _, link := range links {
			amazonLinks[link] = struct{}{}
		}
	}

	saveDir := filepath.Join(dest, time.Now().Format("20060102"))
	retry := uint(0)
	for al := range amazonLinks {
		for {
			res, err := http.Get(al.Href)
			if err != nil {
				retry++
				if retry > 3 {
					return errors.Wrapf(err, "retry count over itbook detail request : %s", al.Text)
				}
				fmt.Println("retry request", al.Text)
				time.Sleep(3 * time.Second)
				continue
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("itbook detail request error: %d", res.StatusCode)
			}

			html, err := toString(res)
			if err != nil {
				return err
			}

			if err = writeFile(saveDir, al.Text+".html", html); err != nil {
				return err
			}
			break
		}
		time.Sleep(1 * time.Second)
	}

	return saveAmazonLinks(saveDir, "link_list.jsonl", amazonLinks)
}

func getLinks(r io.Reader) ([]cmd.AmaoznLink, error) {
	dateStr := time.Now().Format("2006/1/2")
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
		return nil, errors.Wrapf(err, "failed parse response body : %s", string(b))
	}

	var dates []string
	doc.Find("span").Each(func(_ int, s *goquery.Selection) {
		class, exists := s.Attr("class")
		if !exists || class != "a-size-base a-color-secondary a-text-normal" {
			return
		}
		pubDate := strings.TrimSpace(s.Text())
		if pubDate == dateStr {
			dates = append(dates, pubDate)
		}
	})

	if len(dates) == 0 {
		return nil, nil
	}

	var links []cmd.AmaoznLink
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		if len(dates) == len(links) {
			return
		}

		class, exists := s.Attr("class")
		if !exists || class != "a-link-normal a-text-normal" {
			return
		}

		target, exists := s.Attr("target")
		if !exists || target != "_blank" {
			return
		}

		href, exists := s.Attr("href")
		if exists {
			text := strings.TrimSpace(s.Text())
			l := cmd.AmaoznLink{
				Text: text,
				Href: "https://www.amazon.co.jp/" + href,
			}
			links = append(links, l)
		}
	})
	return links, nil
}

func toString(resp *http.Response) (string, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(body)
	html := buf.String()
	return html, nil
}

func writeFile(dir, fileName, html string) error {
	f, err := cmd.CreateFile(filepath.Join(dir, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(html); err != nil {
		return errors.Wrap(err, "failed to write html")
	}

	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}

func saveAmazonLinks(out, fileName string, links map[cmd.AmaoznLink]struct{}) error {
	f, err := cmd.CreateOutFile(filepath.Join(out, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	for link := range links {
		err = cmd.AppendOutFile(f, link)
		if err != nil {
			return err
		}
	}
	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}
