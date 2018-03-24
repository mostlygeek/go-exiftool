package exiftool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStayOpen(t *testing.T) {
	assert := assert.New(t)

	stayopen, err := NewStayopen("exiftool")
	if !assert.NoError(err) {
		return
	}

	meta, err := stayopen.Extract("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}

	create, ok := meta.CreateDate()
	assert.True(ok)
	assert.Equal("2016-06-17 19:16:43 +0000 UTC", create.String())

	stayopen.Stop()

	// extracting after stop should fail
	_, err = stayopen.Extract("testdata/IMG_7238.JPG")
	assert.Error(err)
}

func TestStayOpenErrorsOnBadBin(t *testing.T) {
	_, err := NewStayopen("not.a.rea.bin")
	assert.Error(t, err)
}
