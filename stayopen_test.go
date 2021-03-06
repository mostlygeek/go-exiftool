package exiftool

import (
	"bufio"
	"io"

	"testing"

	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestStayOpen(t *testing.T) {
	assert := assert.New(t)

	stayopen, err := NewStayOpen("exiftool", "-json")
	if !assert.NoError(err) {
		return
	}

	data, err := stayopen.Extract("testdata/IMG_7238.JPG")
	if !assert.NoError(err) {
		return
	}
	createDate, err := jsonparser.GetString(data, "[0]", "CreateDate")
	if assert.NoError(err) {
		assert.Equal("2016:06:17 19:16:43", createDate)
	}

	data, err = stayopen.ExtractFlags("testdata/IMG_7238.JPG", "-ShutterSpeed")
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

	stayopen.Stop()

	// extracting after stop should fail
	_, err = stayopen.Extract("testdata/IMG_7238.JPG")
	assert.Error(err)
}

func TestStayOpenErrorsOnBadBin(t *testing.T) {
	_, err := NewStayOpen("not.a.rea.bin")
	assert.Error(t, err)
}

func TestSplitReadyToken(t *testing.T) {
	assert := assert.New(t)
	data := []byte("xxx{ready}\nyyy\n{ready}\nzzz{ready}\n")

	advance, token, err := splitReadyToken(data, false)
	assert.NoError(err)
	assert.Equal(11, advance)
	assert.Equal([]byte("xxx"), token)

	// note that yyy has extra line ending characters
	data = data[advance:]
	advance, token, err = splitReadyToken(data, false)
	assert.NoError(err)
	assert.Equal(12, advance)
	assert.Equal([]byte("yyy\n"), token)

	data = data[advance:]
	advance, token, err = splitReadyToken(data, true)
	assert.Equal(bufio.ErrFinalToken, err)
	assert.Equal(11, advance)
	assert.Equal([]byte("zzz"), token)
}

func TestSplitReadyTokenWindows(t *testing.T) {
	assert := assert.New(t)

	data := []byte("xxx{ready}\r\nyyy\r\n{ready}\r\nzzz{ready}\r\n")
	advance, token, err := splitReadyToken(data, false)
	assert.NoError(err)
	assert.Equal(12, advance)
	assert.Equal([]byte("xxx"), token)

	// note that yyy has extra line ending characters
	data = data[advance:]
	advance, token, err = splitReadyToken(data, false)
	assert.NoError(err)
	assert.Equal(14, advance)
	assert.Equal([]byte("yyy\r\n"), token)

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
	assert.Equal([]byte("--\n"), token)
	assert.Equal(bufio.ErrFinalToken, err)
}

func TestSplitReadyTokenNoData(t *testing.T) {
	assert := assert.New(t)
	data := []byte("{ready}\n{ready}\n")
	advance, token, err := splitReadyToken(data, false)
	assert.Equal(8, advance)
	assert.Len(token, 0)
	assert.NoError(err)

	advance, token, err = splitReadyToken(data, true)
	assert.Equal(8, advance)
	assert.Len(token, 0)
	assert.NoError(err)
}
