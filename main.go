package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliergopher/grab/v3"
)

/**

Only works on https://ww3.mangakakalot.tv/

**/

func downloadAllImageaPage(url string) {
	// open url
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// create reader
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	c := 1
	doc.Find("#vungdoc img").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("data-src")
		resp, _ := grab.Get(".", link)
		os.Rename(resp.Filename, strconv.Itoa(c)+".jpg")
		c += 1
	})

}

func getTitleAndAllChapterLinkPage(link string) (string, []string) {
	ret := []string{}
	res, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// create reader
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#chapter > div > div.chapter-list a").Each(func(i int, s *goquery.Selection) {
		l := s.AttrOr("href", "")
		ret = append(ret, l)
	})

	return doc.Find("body > div.container > div.main-wrapper > div.leftCol > div.manga-info-top > ul > li:nth-child(1) > h1").First().Text(), ret
}

func main() {
	var link_page string
	outdir := "download"
	flag.StringVar(&link_page, "link", "", "the link to the manga main page")
	flag.StringVar(&outdir, "o", "", "output directory")

	flag.Parse()

	if link_page == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	title, links := getTitleAndAllChapterLinkPage(link_page)
	fmt.Println(links)
	fmt.Println(title)
	os.Chdir(outdir)
	os.Mkdir(title, 0755)
	os.Chdir(title)
	for _, l := range links {
		l2 := strings.Split(l, "-")
		folder_name := l2[len(l2)-1]
		fmt.Println(folder_name)
		os.Mkdir(folder_name, 0755)
		os.Chdir(folder_name)
		u, _ := url.Parse(link_page)
		downloadAllImageaPage("https://" + u.Host + l)
		os.Chdir("..")
	}
}
