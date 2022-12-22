package chain

import (
	pairV2 "dex-tool/src/configs/contract/pair-v2"
	pairV3 "dex-tool/src/configs/contract/pair-v3"
	"dex-tool/src/models/entities"
)

func GetRpc(rpc []string) string {
	return rpc[0]
}

func GetPairAbi(version int) string {
	switch version {
	case 2:
		return pairV2.PairV2MetaData.ABI
	case 3:
		return pairV3.PairV3MetaData.ABI
	default:
		return ""
	}
}

func GetTopic(topics []entities.TopicBlock, version int) string {
	for _, topic := range topics {
		if topic.Version == version {
			return topic.Topic
		}
	}
	return ""
}

func GetCurrentBlock(topics []entities.TopicBlock, version int) int {
	for _, topic := range topics {
		if topic.Version == version {
			return topic.Block
		}
	}
	return -1
}


