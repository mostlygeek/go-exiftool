package exiftool

import (
	"os"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestExtractFlags(t *testing.T) {
	assert := assert.New(t)
	data, err := Extract("exiftool", "testdata/IMG_7238.JPG", "-j", "-CircleOfConfusion")
	if !assert.NoError(err) {
		return
	}
	coc, err := jsonparser.GetString(data, "[0]", "CircleOfConfusion")
	assert.NoError(err)
	assert.Equal("0.004 mm", coc)

}

func testExtractReaderFlags(t *testing.T) {
	assert := assert.New(t)
	f, err := os.Open("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}
	data, err := ExtractReader("exiftool", f, "-j", "-CircleOfConfusion")
	if !assert.NoError(err) {
		return
	}
	coc, err := jsonparser.GetString(data, "[0]", "CircleOfConfusion")
	assert.NoError(err)
	assert.Equal("0.004 mm", coc)
}
