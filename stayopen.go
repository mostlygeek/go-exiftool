package exiftool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
)

// Stayopen abstracts running exiftool with `-stay_open` to greatly improve
// performance by removing the startup overhead of perl for each indivdual file
// remember to call Stayopen.Stop() to signal exiftool to shutdown otherwise
// it may make zombie perl processes
type Stayopen struct {
	sync.Mutex
	cmd     *exec.Cmd
	waitEnd sync.WaitGroup
	in      chan string
	out     chan []byte
}

// Extract calls exiftool on the supplied filename
func (e *Stayopen) Extract(filename string) (*Metadata, error) {
	e.Lock()
	defer e.Unlock()

	if e.cmd == nil {
		return nil, errors.New("Stopped")
	}

	// send it and wait for it to come back from exiftool
	e.in <- filename
	data := <-e.out

	container := make([]map[string]interface{}, 1, 1)
	err := json.Unmarshal(data, &container)
	if err != nil {
		return nil, errors.Wrap(err, "JSON unmarshal failed")
	}

	if len(container) != 1 {
		return nil, errors.New("Expected one record")
	}

	meta := NewMetadata(container[0])
	if errstr := meta.Error(); errstr != "" {
		return meta, errors.New(errstr)
	}

	return meta, nil
}

func (e *Stayopen) Stop() {
	e.Lock()
	defer e.Unlock()

	// closing the in channel will trigger a shutdown
	// wait for both goroutines to finish before finishing
	close(e.in)
	e.waitEnd.Wait()
	e.cmd = nil
}

func NewStayopen(exiftool string) *Stayopen {
	stayopen := &Stayopen{
		in:  make(chan string),
		out: make(chan []byte),
	}

	stayopen.cmd = exec.Command(exiftool, "-stay_open", "True", "-@", "-")
	stdin, _ := stayopen.cmd.StdinPipe()
	stdout, _ := stayopen.cmd.StdoutPipe()

	var startReady sync.WaitGroup
	startReady.Add(2)

	// send commands to exiftool's stdin
	go func() {
		stayopen.cmd.Start()
		startReady.Done()
		stayopen.waitEnd.Add(1)
		for filename := range stayopen.in {
			fmt.Fprintln(stdin, "-json")
			fmt.Fprintln(stdin, "--printConv")
			fmt.Fprintln(stdin, filename)
			fmt.Fprintln(stdin, "-execute")
		}

		// write message telling it to close
		// but don't actually wait for the command to stop
		fmt.Fprintln(stdin, "-stay_open")
		fmt.Fprintln(stdin, "False")
		fmt.Fprintln(stdin, "-execute")

		// closing stdout will stop the scanner goroutine
		stdout.Close()
		stayopen.waitEnd.Done()
	}()

	// scan exiftool's stdout, parse out JSON messages
	// and publish them on the out channel
	go func() {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(splitReadyToken)

		startReady.Done()
		stayopen.waitEnd.Add(1)

		for scanner.Scan() {
			results := scanner.Bytes()
			sendResults := make([]byte, len(results), len(results))
			copy(sendResults, results)
			stayopen.out <- sendResults
		}
		close(stayopen.out)
		stayopen.waitEnd.Done()
	}()

	// wait for both go-routines to startup
	startReady.Wait()
	return stayopen
}

func splitReadyToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.Index(data, []byte("\n{ready}")); i >= 0 {
		return i + 8, data[:i], nil
	}

	if atEOF {
		return len(data), data, io.EOF
	}

	return 0, nil, nil
}
