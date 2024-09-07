package db

import (
	"time"

	"gorm.io/gorm"
)

type SearchIndex struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4()"`
	Value     string
	Urls      []CrawledUrl   `gorm:"many2many:token_urls"`
	CreatedAt *time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (s *SearchIndex) TableName() string {
	return "search_index"
}

func (s *SearchIndex) Save(index map[string][]string, crawledUrls []CrawledUrl) error {
	for value, ids := range index {
		newIndex := &SearchIndex{Value: value}
		if err := DBConn.Where(SearchIndex{Value: value}).Find(newIndex).Error; err != nil {
			return err		
		}

		var urlsToAppend []CrawledUrl
		for _, id := range ids {
			for _, url := range crawledUrls {
				if url.ID == id {
					urlsToAppend = append(urlsToAppend, url)
					break
				}
			}
		}
		if err := DBConn.Model(&newIndex).Association("Urls").Append(&urlsToAppend); err != nil {
			return err
		}
	} 
	return nil
}
