package main

import (
	"dex-tool/src/configs/db"
	"dex-tool/src/configs/env"
	"dex-tool/src/crawler"
	"dex-tool/src/models/repositories"
	"dex-tool/src/task"
	_chain "dex-tool/src/utils/chain"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	redisConnection := asynq.RedisClientOpt{
		Addr:     env.ConfigEnv.Redis.Address,
		Username: env.ConfigEnv.Redis.Username,
		Password: env.ConfigEnv.Redis.Password,
		DB:       env.ConfigEnv.Redis.Db,
	}
	asynqClient := asynq.NewClient(redisConnection)
	defer func(client *asynq.Client) {
		err := client.Close()
		if err != nil {

		}
	}(asynqClient)

	redisInstance := db.GetRedisInstance()

	argsWithProg := os.Args
	db.Setup(env.ConfigEnv.Db.Url, env.ConfigEnv.Db.Name)
	log.Println(argsWithProg[1])
	chainRepo := repositories.NewChainRepository()


	var blockPerProcess = env.ConfigEnv.Crawl.MaxBlock
	var safeBlock = env.ConfigEnv.Crawl.SafeBlock
	for {
		chain, err := chainRepo.GetChainByName(argsWithProg[1])
		if err != nil {
			panic(err)
		}
		topicBlock := chain.TopicBlock

		block, err := redisInstance.Get(crawler.GetKeyBlock(chain.Name))
		if err != nil {
			panic(err)
		}
		blockNumber, _ := strconv.Atoi(block)
		blockNumber = blockNumber - safeBlock
		currentBlockV2 := _chain.GetCurrentBlock(topicBlock, 2)
		currentBlockV3 := _chain.GetCurrentBlock(topicBlock, 3)
		var toBlockV2, toBlockV3 int
		if currentBlockV2 + blockPerProcess > blockNumber {
			toBlockV2 = blockNumber
		} else {
			toBlockV2 = currentBlockV2 + blockPerProcess
		}

		if currentBlockV3 + blockPerProcess > blockNumber {
			toBlockV3 = blockNumber
		} else {
			toBlockV3 = currentBlockV3 + blockPerProcess
		}

		taskV2 := task.NewCrawlTransactionTask(currentBlockV2, toBlockV2)
		taskV3 := task.NewCrawlTransactionTask(currentBlockV3, toBlockV3)

		if _, err := asynqClient.Enqueue(taskV2, asynq.Queue(fmt.Sprintf("%s_%s", chain.Name, _chain.GetTopic(topicBlock, 2)))); err != nil {
			panic(err)
		}

		if _, err := asynqClient.Enqueue(taskV3, asynq.Queue(fmt.Sprintf("%s_%s", chain.Name, _chain.GetTopic(topicBlock, 3)))); err != nil {
			panic(err)
		}

		chainRepo.UpdateBlock(chain.Name, 2, toBlockV2)
		chainRepo.UpdateBlock(chain.Name, 3, toBlockV3)
		time.Sleep(3 * time.Second)
	}
}
