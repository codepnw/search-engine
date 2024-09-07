package search

import (
	"fmt"
	"time"

	"github.com/codepnw/search-engine/internal/db"
)

func RunEngine() {
	fmt.Println("started search engine crawl...")
	defer fmt.Println("search engine crawl has finished")

	settings := &db.SearchSettings{}

	if err := settings.Get(); err != nil {
		fmt.Println("something went wrong getting settings")
		return
	}

	if !settings.SearchOn {
		fmt.Println("search is turned off")
		return
	}

	crawl := &db.CrawledUrl{}
	nextUrls, err := crawl.GetNextCawlUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("something went wrong get next urls")
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()
	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdatedUrl(db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         false,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Heading:         result.CrawlData.Heading,
				LastTested:      &testedTime,
			})
			if err != nil {
				fmt.Println("something went wrong updating failed url")
			}
			continue
		}
		// Success
		err := next.UpdatedUrl(db.CrawledUrl{
			ID:              next.ID,
			Url:             next.Url,
			Success:         result.Success,
			CrawlDuration:   result.CrawlData.CrawlTime,
			ResponseCode:    result.ResponseCode,
			PageTitle:       result.CrawlData.PageTitle,
			PageDescription: result.CrawlData.PageDescription,
			Heading:         result.CrawlData.Heading,
			LastTested:      &testedTime,
		})
		if err != nil {
			fmt.Println("something went wrong updating success url")
			fmt.Println(next.Url)
		}
		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	} // end of range
	if !settings.AddNew {
		return
	}

	// Insert new urls
	for _, newUrl := range newUrls {
		if err := newUrl.Save(); err != nil {
			fmt.Println("something went wrong adding new url to database")
		}
	}
	fmt.Printf("\n Added %d new urls to database", len(newUrls))
}

func RunIndex() {
	fmt.Println("started search indexing...")
	defer fmt.Println("search indexing has finished")

	crawled := &db.CrawledUrl{}
	notIndexed, err := crawled.GetNotIndex()
	if err != nil {
		return
	}

	idx := make(Index)
	idx.Add(notIndexed)
	searchIndex := &db.SearchIndex{}

	if err := searchIndex.Save(idx, notIndexed); err != nil { 
		return
	}

	if err := crawled.SetIndexedTrue(notIndexed); err != nil {
		return
	}
}