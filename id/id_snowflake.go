package id

import (
	"os"

	"gdcx.com/infra/logger"
	"github.com/bwmarrin/snowflake"
)

type snowFlakeID struct {
	node *snowflake.Node
}

func newSnowFlakeID(node int64) *snowFlakeID {
	if node == 0 {
		host, _ := os.Hostname()
		node = str2int(host, -1^(-1<<snowflake.NodeBits))
	}

	n, err := snowflake.NewNode(node)
	if err != nil {
		logger.Fatalf("can't create id generator(%v)", err)
	}
	return &snowFlakeID{
		node: n,
	}
}

func (s *snowFlakeID) NextID() int64 {
	return int64(s.node.Generate())
}

func (s *snowFlakeID) NextIDStr() string {
	return s.node.Generate().Base36()
}

// compute hashcode of string
func str2int(str string, maxint int64) int64 {
	var sum int64
	for i := 0; i < len(str); i++ {
		sum = 31*sum + int64(str[i])
	}
	sum %= maxint
	if sum < 0 {
		return ^sum + 1
	}
	return sum
}
