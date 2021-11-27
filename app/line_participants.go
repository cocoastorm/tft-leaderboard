package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type LineParticipantsPool struct {
	lines io.ReadCloser
	err error
}

func NewParticipantsFile(filename string) *LineParticipantsPool {
	pool := LineParticipantsPool{}

	if f, err := os.Open(filename); err == nil {
		pool.lines = f
	} else {
		pool.err = fmt.Errorf("failed to open participants file located at: %s, %s", filename, err)
	}
	
	return &pool
}

func (pool *LineParticipantsPool) GetParticipants() []string {
	names := make([]string, 0)
	defer pool.lines.Close()

	scanner := bufio.NewScanner(pool.lines)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			names = append(names, text)
		}
	}

	return names
}

func (lines *LineParticipantsPool) Error() error {
	return lines.err
}
