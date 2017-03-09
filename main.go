package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s URL", os.Args[0])
	}
	response, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		z := html.NewTokenizer(response.Body)

		for {
			tt := z.Next()

			switch {
			case tt == html.ErrorToken:
				// End of the document, we're done
				return
			case tt == html.StartTagToken:
				t := z.Token()

				if t.Data == "a" {
					for _, attr := range t.Attr {
						if attr.Key == "href" {
							log.Printf("%+v", attr.Val)
						}
					}
				}
			}
		}
	}
}
