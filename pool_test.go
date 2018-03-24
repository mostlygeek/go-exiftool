package exiftool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)

	pool, err := NewPool("exiftool", 2)
	if !assert.NoError(err) {
		return
	}

	meta, err := pool.Extract("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}

	create, ok := meta.CreateDate()
	assert.True(ok)
	assert.Equal("2016-06-17 19:16:43 +0000 UTC", create.String())

	pool.Stop()

	// extracting after stop should fail
	_, err = pool.Extract("testdata/IMG_7238.JPG")
	assert.Error(err)
}

func TestPoolErrorsOnBadBin(t *testing.T) {
	_, err := NewPool("not.a.rea.bin", 1)
	assert.Error(t, err)
}
