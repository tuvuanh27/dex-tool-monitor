package task

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
)

func HandleCrawlTransactionTask(c context.Context, t *asynq.Task) error {
	payload := t.Payload()
	var crawlTransaction CrawlTransaction
	err := json.Unmarshal(payload, &crawlTransaction)
	if err != nil {
		return err
	}
	log.Println("Crawl Transaction Task Payload", crawlTransaction)

	return nil
}
