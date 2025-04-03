package deployer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
)

func cleanupLine(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\t", " ")
}

type LoggerHandler struct {
	lines []string

	LogCallback func(string)
}

func (h *LoggerHandler) IngestReaders(done chan<- struct{}, stdout io.Reader, stderr io.Reader) error {
	var wg sync.WaitGroup
	wg.Add(2)

	ingest := make(chan string)

	engine := func(r io.Reader) {
		s := bufio.NewScanner(r)

		for s.Scan() {
			txt := s.Text()
			h.LogCallback(cleanupLine(txt))
			ingest <- txt

		}
		wg.Done()
	}

	go engine(stdout)
	go engine(stderr)

	go func() {
		wg.Wait()
		close(ingest)
	}()

	go func() {
		for line := range ingest {
			h.lines = append(h.lines, line)
		}
		close(done)
	}()

	return nil
}

func (h *LoggerHandler) AugmentError(err error) error {
	return fmt.Errorf("\n%s\n%w", strings.Join(h.lines, "\n"), err)
}
