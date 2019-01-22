package exiftool

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// Stayopen abstracts running exiftool with `-stay_open` to greatly improve
// performance. Remember to call Stayopen.Stop() to signal exiftool to shutdown
// to avoid zombie perl processes
type Stayopen struct {
	l   sync.Mutex
	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser

	// default flags to pass to every extract call
	defaultFlags string

	scanner *bufio.Scanner
}

// Extract calls exiftool on the supplied filename
func (e *Stayopen) Extract(filename string) ([]byte, error) {
	return e.ExtractFlags(filename)
}

func (e *Stayopen) ExtractFlags(filename string, flags ...string) ([]byte, error) {
	e.l.Lock()
	defer e.l.Unlock()

	if e.cmd == nil {
		return nil, errors.New("Stopped")
	}

	if !strconv.CanBackquote(filename) {
		return nil, ErrFilenameInvalid
	}

	// send the request
	fmt.Fprintln(e.stdin, e.defaultFlags)
	if len(flags) > 0 {
		fmt.Fprintln(e.stdin, strings.Join(flags, "\n"))
	}
	fmt.Fprintln(e.stdin, filename)
	fmt.Fprintln(e.stdin, "-execute")

	if !e.scanner.Scan() {
		return nil, errors.New("Failed to read output")
	} else {
		results := e.scanner.Bytes()
		sendResults := make([]byte, len(results), len(results))
		copy(sendResults, results)
		return sendResults, nil
	}

}

func (e *Stayopen) Stop() {
	e.l.Lock()
	defer e.l.Unlock()

	// write message telling it to close
	// but don't actually wait for the command to stop
	fmt.Fprintln(e.stdin, "-stay_open")
	fmt.Fprintln(e.stdin, "False")
	fmt.Fprintln(e.stdin, "-execute")
	e.cmd = nil
}

func NewStayOpen(exiftool string, flags ...string) (*Stayopen, error) {

	var defaultFlags string
	if len(flags) > 0 {
		defaultFlags = strings.Join(flags, "\n")
	}
	stayopen := &Stayopen{
		defaultFlags: defaultFlags,
	}

	stayopen.cmd = exec.Command(exiftool, "-stay_open", "True", "-@", "-")

	stdin, err := stayopen.cmd.StdinPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting stdin pipe")
	}

	stdout, err := stayopen.cmd.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "Failed getting stdout pipe")
	}

	stayopen.stdin = stdin
	stayopen.stdout = stdout
	stayopen.scanner = bufio.NewScanner(stdout)
	stayopen.scanner.Split(splitReadyToken)

	if err := stayopen.cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "Failed starting exiftool in stay_open mode")
	}

	// wait for both go-routines to startup
	return stayopen, nil
}

func splitReadyToken(data []byte, atEOF bool) (int, []byte, error) {
	if tokenPos := bytes.Index(data, []byte("{ready}")); tokenPos >= 0 {

		// incomplete data, no line terminator for {ready} token
		if len(data) == tokenPos+7 {
			return 0, nil, nil
		}

		// On Windows the line endings from Perl will be \r\n, whereas
		// on posix systems it should be \n.  If {ready} is followed by
		// a \r, then assume the full token we are looking for is:
		//
		//   \r\n{ready}\r\n
		//
		// otherwise assume it is:
		//
		//   \n{ready}\n
		var tokenSize int
		if len(data) > tokenPos-2 && data[tokenPos-2] == byte('\r') {
			tokenSize = 11          // \r\n{ready}\r\n = 11 bytes
			tokenPos = tokenPos - 2 // strip \r\n from data
		} else {
			// assume \n{ready}\n as the token + line endings
			tokenSize = 9
			tokenPos = tokenPos - 1
		}

		if atEOF && len(data) == (tokenPos+tokenSize) { // nothing left to scan
			return tokenPos + tokenSize, data[:tokenPos], bufio.ErrFinalToken
		} else {
			return tokenPos + tokenSize, data[:tokenPos], nil
		}
	}

	if atEOF {
		return 0, data, io.EOF
	}

	return 0, nil, nil
}
