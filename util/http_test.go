package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadSum(t *testing.T) {
	assert := assert.New(t)

	file := "https://godouav.oss-cn-shenzhen.aliyuncs.com/gproxy/0.1.2/linux/amd64/edge"
	sha256 := "5aebfb9442fb310f739ae9d948e6b93131549356d24785471a5b1138a446b65e"

	sum, err := HttpDownload(file, nil, "/tmp/out", Sha256)
	assert.NoError(err)
	assert.Equal(sum, sha256)
}
