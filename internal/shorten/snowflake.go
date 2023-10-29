package shorten

import (
	"github.com/bwmarrin/snowflake"
	"os"
	"strconv"
	"time"
)

func getNewNode() (*snowflake.Node, error) {
	nodeId, found := os.LookupEnv("NODE_ID")
	if !found {
		nodeId = "1"
	}

	id, err := strconv.ParseInt(nodeId, 10, 64)
	if err != nil {
		return nil, err
	}

	return snowflake.NewNode(id)
}

// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+
// | 0           | 0 ... 10011010100111010101011000 | 00000001 | 000000000000 |
// +--------------------------------------------------------------------------+
func GenerateSnowFlakeKey() (int64, error) {
	snowflake.Epoch = time.Date(2023, 10, 29, 0, 0, 0, 0, time.UTC).UnixMilli()
	node, err := getNewNode()
	if err != nil {
		return 0, err
	}
	id := node.Generate()
	return int64(id), nil
}
