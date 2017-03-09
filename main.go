package main

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s URL", os.Args[0])
	}

	// Create channels for message passing.
	messages := make(chan string)
	finished := make(chan bool)
	workers := 10

	// Pass in init url
	messages <- os.Args[1]

	for i := 0; i < workers; {
		go func() {
			defer func() {
				finished <- true
			}()
			inc := <-messages
			log.Printf("%+v", inc)

			urls, err := ScrapeUrl(inc)
			if err != nil {
				log.Fatal(err)
			}
			for _, u := range urls {
				messages <- u
			}
		}()

		select {
		case u := <-messages:
			log.Printf("Recieved: %+v", u)
		case <-finished:
			i++
		}
	}
}

func ScrapeUrl(uri string) ([]string, error) {
	response, err := http.Get(uri)
	ret := []string{}

	if err != nil {
		return nil, err
	} else {
		defer response.Body.Close()
		z := html.NewTokenizer(response.Body)

		for {
			tt := z.Next()

			switch {
			case tt == html.ErrorToken:
				// End of the document, we're done
				return ret, nil
			case tt == html.StartTagToken:
				t := z.Token()

				if t.Data == "a" {
					for _, attr := range t.Attr {
						if attr.Key == "href" {
							u, err := url.ParseRequestURI(attr.Val)
							if err != nil {
								continue
							} else {
								if u.IsAbs() {
									log.Printf("Found %+v", attr.Val)
									ret = append(ret, attr.Val)
								}
							}
						}
					}
				}
			}
		}
	}

	return ret, nil
}
