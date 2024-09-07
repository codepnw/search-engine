package utils

import (
	"fmt"

	"github.com/codepnw/search-engine/internal/search"
	"github.com/robfig/cron"
)

func StartCronJobs() {
	c := cron.New()
	c.AddFunc("0 * * * *", search.RunEngine) // Run Every Hour
	c.AddFunc("15 * * * *", search.RunIndex) // Run Every Hour at 15 mins past
	c.Start()
	cronCount := len(c.Entries())
	fmt.Printf("setup %d cron jobs\n", cronCount)
}