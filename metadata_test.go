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

	{ // with gps coords, timezone info is correct
		metaGeo, _ := Extract("testdata/IMG_7238.JPG")
		create, ok := metaGeo.CreateDate()
		assert.True(ok)
		assert.Equal("2016-06-17 19:16:43 +0100 BST", create.String())
	}

	{ // no gps coords get +000 UTC timezone
		metaNoGeo, _ := Extract("testdata/IMG_7238-nogeo.jpg")
		create, ok := metaNoGeo.CreateDate()
		assert.True(ok)
		assert.Equal("2016-06-17 19:16:43 +0000 UTC", create.String())
	}
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
	assert.Equal(51.49935555555555, lat)
	assert.Equal(-0.12998055555555554, long)
}

func TestParseGPS(t *testing.T) {
	c := map[string]float64{
		`51 deg 29' 57.68" N`: 51.49935555,
		`51 deg 29' 57.68" S`: -51.4993555,
		`10 deg 20" S`:        -10.0055556,
		`10 deg 20"`:          10.0055556,
		`10 deg`:              10.0,
		`0 deg 20'`:           0.3333333,
		`0 deg 20' S`:         -0.3333333,
	}

	assert := assert.New(t)
	for coord, expect := range c {
		val, err := ParseGPS(coord)
		assert.NoError(err)
		assert.InDelta(expect, val, 0.0000001)
	}
}
