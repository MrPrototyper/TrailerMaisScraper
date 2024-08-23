package main

import (
	"fmt"
	"log"
	"net/http"
	"scrap-data/db"
	"scrap-data/movie"
	"scrap-data/movie/models"
	"strconv"

	"golang.org/x/net/html"
)

func fetchHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page: status code %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func extractMaxPage(n *html.Node) int {
	if n.Type == html.ElementNode && n.Data == "li" {
		// Check if the li element contains a link with class "paginate-page"
		for _, attr := range n.Attr {
			if attr.Key == "class" && attr.Val == "paginate-page" {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.ElementNode && c.Data == "a" {
						pageNum, err := strconv.Atoi(c.FirstChild.Data)
						if err == nil {
							return pageNum
						}
					}
				}
			}
		}
	}

	// Recursively search through child nodes
	maxPage := 0
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		page := extractMaxPage(c)
		if page > maxPage {
			maxPage = page
		}
	}
	return maxPage
}

func extractTitles(n *html.Node) {
	client, ctx := db.GetClient()

	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "alt" {
				fmt.Println(attr.Val)
				newMovie := models.Movie{
					Title: attr.Val,
				}
				id, err := movie.Save(client, ctx, newMovie)
				if err != nil {
					log.Fatalf("Failed to save user: %v", err)
				}
				fmt.Printf("User saved with ID: %v\n", id)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTitles(c)
	}
}

func getAllPagesByYear(year int) {
	url := fmt.Sprintf("https://letterboxd.com/sprudelheinz/list/all-the-movies-sorted-by-movie-posters-1/year/%v/page/1/", year)
	firstPage, firstPageErr := fetchHTML(url)

	if firstPageErr != nil {
		log.Fatalf("Error fetching HTML: %v", firstPageErr)
	}

	maxNumPage := extractMaxPage(firstPage)

	fmt.Printf("the year %d has %d pages", year, maxNumPage)

	for i := 1; i <= maxNumPage; i++ {
		url := fmt.Sprintf("https://letterboxd.com/sprudelheinz/list/all-the-movies-sorted-by-movie-posters-1/year/%v/page/%d/", year, i)
		doc, err := fetchHTML(url)

		if err != nil {
			log.Fatalf("Error fetching HTML: %v", err)
		}

		extractTitles(doc)
	}
}

func main() {
	_, _, _, connErr := db.Connect()

	if connErr != nil {
		log.Fatalf("Error Connecting to DB: %v", connErr)
	}

	startYear := 2023
	endYear := 2024

	for year := startYear; year <= endYear; year++ {
		getAllPagesByYear(year)
	}
}
