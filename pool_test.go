package exiftool

import (
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	assert := assert.New(t)

	pool, err := NewPool("exiftool", 2, "-json")
	if !assert.NoError(err) {
		return
	}

	data, err := pool.Extract("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}

	createDate, err := jsonparser.GetString(data, "[0]", "CreateDate")
	if assert.NoError(err) {
		assert.Equal("2016:06:17 19:16:43", createDate)
	}

	data, err = pool.ExtractFlags("testdata/IMG_7238.JPG", "-ShutterSpeed")
	if !assert.NoError(err) {
		return
	}

	if ss, err := jsonparser.GetString(data, "[0]", "ShutterSpeed"); assert.NoError(err) {
		assert.Equal("1/123", ss)
	}

	// make sure there's nothing in data other than shutterspeed
	createDate, err = jsonparser.GetString(data, "[0]", "CreateDate")
	assert.NotNil(err)
	assert.Equal("", createDate)

	pool.Stop()

	// extracting after stop should fail
	_, err = pool.Extract("testdata/IMG_7238.JPG")
	assert.Error(err)
}

func TestPoolErrorsOnBadBin(t *testing.T) {
	_, err := NewPool("not.a.rea.bin", 1)
	assert.Error(t, err)
}
