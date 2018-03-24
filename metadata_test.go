package exiftool

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	meta *Metadata
)

func init() {
	// reuse to improve test speed
	var err error
	meta, err = Extract("testdata/IMG_7238.JPG")
	if err != nil {
		panic(err.Error())
	}

}

func TestMIMEType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("image/jpeg", meta.MIMEType())
}

func TestCreateDate(t *testing.T) {
	assert := assert.New(t)
	create, ok := meta.CreateDate()
	assert.True(ok)
	assert.Equal("2016-06-17 19:16:43 +0000 UTC", create.String())
}

func TestGetBytes(t *testing.T) {
	assert := assert.New(t)
	thumb, err := meta.GetBytes("EXIF", "ThumbnailImage")
	assert.NoError(err)
	assert.Len(thumb, 9898)
	sha := sha1.Sum(thumb)
	assert.Equal("3b06e0f201303721a866c8f816bf889f99adcb6d", hex.EncodeToString(sha[:]))
}

func TestGPSPosition(t *testing.T) {
	assert := assert.New(t)

	lat, long, ok := meta.GPSPosition()
	assert.True(ok)
	assert.Equal(51.4993555555556, lat)
	assert.Equal(-0.129980555555556, long)
}
