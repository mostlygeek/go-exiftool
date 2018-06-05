package exiftool

import (
	"bufio"
	"io"
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
	assert.Equal("2016-06-17 19:16:43 +0100 BST", create.String())

	stayopen.Stop()

	// extracting after stop should fail
	_, err = stayopen.Extract("testdata/IMG_7238.JPG")
	assert.Error(err)
}

func TestStayOpenErrorsOnBadBin(t *testing.T) {
	_, err := NewStayopen("not.a.rea.bin")
	assert.Error(t, err)
}

func TestSplitReadyToken(t *testing.T) {
	assert := assert.New(t)
	data := []byte("xxx\n{ready}\nyyy\n{ready}\nzzz\n{ready}\n")

	advance, token, err := splitReadyToken(data, false)
	assert.NoError(err)
	assert.Equal(12, advance)
	assert.Equal([]byte("xxx"), token)

	data = data[advance:]
	advance, token, err = splitReadyToken(data, false)
	assert.NoError(err)
	assert.Equal(12, advance)
	assert.Equal([]byte("yyy"), token)

	data = data[advance:]
	advance, token, err = splitReadyToken(data, true)
	assert.Equal(bufio.ErrFinalToken, err)
	assert.Equal(12, advance)
	assert.Equal([]byte("zzz"), token)
}

// TestSplitReadyTokenPartial tests that more data is requested
// when we don't have a full delimter yet
func TestSplitReadyTokenPartial(t *testing.T) {
	assert := assert.New(t)
	data := []byte("xxx\n{ready}")                      // missing \n
	advance, token, err := splitReadyToken(data, false) // not at EOF
	assert.Equal(0, advance)
	assert.Nil(token)
	assert.Nil(err)
}

// TestSplitReadyToken tests behaviour when we've hit
// EOF on a Reader but we don't have a full delimiter yet
func TestSplitReadyTokenEOF(t *testing.T) {
	assert := assert.New(t)
	data := []byte("xxx\n{ready") // no full token
	advance, token, err := splitReadyToken(data, true)

	// we should get an io.EOF error and get back all
	// the data that couldn't be parsed
	assert.Equal(0, advance)
	assert.Equal(io.EOF, err)
	assert.Equal(data, token)
}

func TestSplitReadyTokenFinalToken(t *testing.T) {
	assert := assert.New(t)
	data := []byte("--\n{ready}\n") // just a ready token
	advance, token, err := splitReadyToken(data, true)
	assert.Equal(11, advance)
	assert.Equal([]byte("--"), token)
	assert.Equal(bufio.ErrFinalToken, err)
}
