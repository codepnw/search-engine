package db

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Url             string         `json:"url" gorm:"unique;not null"`
	Success         bool           `json:"success" gorm:"not null"`
	CrawlDuration   time.Duration  `json:"crawlDuration"`
	ResponseCode    int            `json:"responseCode"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Heading         string         `json:"heading"`
	LastTested      *time.Time     `json:"lastTested"`
	Indexed         bool           `json:"indexed" gorm:"default:false"`
	CreatedAt       *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (c *CrawledUrl) UpdatedUrl(input CrawledUrl) error {
	query := `"url", "success", "crawl_duration", "response_code", "page_title", "page_description", "heading", "last_tested", "updated_at"`

	tx := DBConn.Select(query).Omit("created_at").Save(&input)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return tx.Error
	}
	return nil
}

func (c *CrawledUrl) GetNextCawlUrls(limit int) ([]CrawledUrl, error) {
	var urls []CrawledUrl

	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

func (c *CrawledUrl) Save() error {
	tx := DBConn.Save(&c)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (c *CrawledUrl) GetNotIndex() ([]CrawledUrl, error) {
	var urls []CrawledUrl

	tx := DBConn.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		return []CrawledUrl{}, tx.Error
	}
	return urls, nil
}

func (c *CrawledUrl) SetIndexedTrue(urls []CrawledUrl) error {	
	for _, url := range urls {
		url.Indexed = true
		tx := DBConn.Save(&url)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
