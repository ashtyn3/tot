package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type row struct {
	Id     int
	Text   string
	Length string
	Title  struct {
		Author,
		Title string
		Type string
	}
}

func title(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc.Find("p").Next().First().Text()
}

func main() {
	res, err := http.Get("https://web.archive.org/web/20210115204923/http://typeracerdata.com/texts?texts=full&sort=relative_average")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	items := []row{}
	doc.Find(".stats > tbody").Children().Each(func(n int, s *goquery.Selection) {
		if n == 0 {
			return
		}

		url, _ := s.Children().Find("a").First().Attr("href")
		from := title("https://web.archive.org/" + url)
		r := regexp.MustCompile(`\-\s*from\s(.*)\,(.*)\s*by(.*)`)
		subs := r.FindAllStringSubmatch(from, 1)
		if len(subs) == 0 {
			return
		} else {
			current := row{}
			id, _ := strconv.Atoi(strings.Replace(s.Children().First().Text(), ".", "", -1))
			current.Id = id
			current.Text = s.Children().Find("a").First().Text()

			current.Title.Title = strings.TrimSpace(subs[0][1])
			current.Title.Type = strings.TrimSpace(strings.TrimPrefix(subs[0][2], "a "))
			current.Title.Author = strings.TrimSpace(subs[0][3])

			current.Length = s.Children().Find("a").Parent().Next().First().Text()
			items = append(items, current)
		}

	})
	bjson, _ := json.Marshal(items)

	os.WriteFile("scraped.json", bjson, 0644)
}
