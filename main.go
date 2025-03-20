package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	whist, err := fetch(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	// Sort
	keys := make([]string, 0, len(whist))
	for key := range whist {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return whist[keys[i]] > whist[keys[j]] })

	// Print
	for _, k := range keys {
		fmt.Printf("%s : %d\n", k, whist[k])
	}
}

func fetch(url string) (map[string]int, error) {
	whist := make(map[string]int)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	t := z.Token()
loopDom:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break loopDom // End of the document,  done
		case tt == html.StartTagToken:
			t = z.Token()
		case tt == html.TextToken:
			if t.Data == "script" {
				continue
			}
			txt := strings.TrimSpace(html.UnescapeString(string(z.Text())))
			wrds := strings.Split(txt, " ")
			for _, w := range wrds {
				whist[w] += 1
			}
		}
	}
	return whist, nil
}
