package exiftool

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
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
	for _, f := range flags {
		fmt.Fprintln(e.stdin, f)
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

	flags = append([]string{"-stay_open", "True", "-@", "-", "-common_args"}, flags...)

	stayopen := &Stayopen{}
	stayopen.cmd = exec.Command(exiftool, flags...)

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
	delimPos := bytes.Index(data, []byte("{ready}\n"))
	delimSize := 8

	// maybe we are on Windows?
	if delimPos == -1 {
		delimPos = bytes.Index(data, []byte("{ready}\r\n"))
		delimSize = 9
	}

	if delimPos == -1 { // still no token found
		if atEOF {
			return 0, data, io.EOF
		} else {
			return 0, nil, nil
		}
	} else {
		if atEOF && len(data) == (delimPos+delimSize) { // nothing left to scan
			return delimPos + delimSize, data[:delimPos], bufio.ErrFinalToken
		} else {
			return delimPos + delimSize, data[:delimPos], nil
		}
	}
}
