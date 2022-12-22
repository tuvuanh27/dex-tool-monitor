package task

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"time"
)

const (
	TypeCrawlTransaction = "crawl:transaction"
)

type CrawlTransaction struct {
	FromBLock int `json:"from_block"`
	ToBlock   int `json:"to_block"`
}

func NewCrawlTransactionTask(fromBlock, toBlock int) *asynq.Task {
	taskCrawlTransaction := CrawlTransaction{fromBlock, toBlock}
	payload, err := json.Marshal(taskCrawlTransaction)
	if err != nil {
		return nil
	}

	return asynq.NewTask(TypeCrawlTransaction, payload, asynq.Retention(24*time.Hour))
}
