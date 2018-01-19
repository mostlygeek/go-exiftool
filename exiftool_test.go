package exiftool

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	meta, err := Extract("testdata/Apple_iPhone_6plus.jpg")
	if assert.NoError(err) && assert.NotNil(meta) {
		create, _ := meta.CreateDate()
		assert.Equal("2014-10-31 13:32:37 +0000 UTC", create.String())
	}
}

func TestExtractReader(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	f, err := os.Open("testdata/Apple_iPhone_6plus.jpg")
	if !assert.NoError(err) {
		return
	}

	meta, err := ExtractReader(f)
	if assert.NoError(err) && assert.NotNil(meta) {
		create, _ := meta.CreateDate()
		assert.Equal("2014-10-31 13:32:37 +0000 UTC", create.String())
	}
}
