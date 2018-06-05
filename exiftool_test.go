package exiftool

import (
	"os"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	meta, err := Extract("testdata/IMG_7238.JPG")
	if assert.NoError(err) && assert.NotNil(meta) {
		create, _ := meta.CreateDate()
		assert.Equal("2016-06-17 19:16:43 +0100 BST", create.String())
	}
}

func TestExtractFlags(t *testing.T) {
	assert := assert.New(t)
	data, err := ExtractFlags("exiftool", "testdata/IMG_7238.JPG", "-j", "-CircleOfConfusion")
	if !assert.NoError(err) {
		return
	}
	coc, err := jsonparser.GetString(data, "[0]", "CircleOfConfusion")
	assert.NoError(err)
	assert.Equal("0.004 mm", coc)

}

func TestExtractReader(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	f, err := os.Open("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}

	meta, err := ExtractReader(f)
	if assert.NoError(err) && assert.NotNil(meta) {
		create, _ := meta.CreateDate()
		assert.Equal("2016-06-17 19:16:43 +0100 BST", create.String())
	}
}

func testExtractReaderFlags(t *testing.T) {
	assert := assert.New(t)
	f, err := os.Open("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}
	data, err := ExtractReaderFlags("exiftool", f, "-j", "-CircleOfConfusion")
	if !assert.NoError(err) {
		return
	}
	coc, err := jsonparser.GetString(data, "[0]", "CircleOfConfusion")
	assert.NoError(err)
	assert.Equal("0.004 mm", coc)
}
