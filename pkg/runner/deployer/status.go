package deployer

import (
	"io"
	"time"
)

type ProgressStatus struct {
	Path     string
	Reader   io.Reader
	Writer   io.Writer
	Size     int
	Callback ProgressStatusCallback

	StartTime   time.Time
	bytesCopied int
}

type ProgressStatusCallback func(path string, copied int, size int, startTime time.Time)

func NewProgressStatus(path string, reader io.Reader, callback ProgressStatusCallback) (*ProgressStatus, error) {
	fsize, err := getReaderSize(reader)

	if err != nil {
		return nil, err
	}

	p := &ProgressStatus{
		Path:      path,
		Reader:    reader,
		Callback:  callback,
		Size:      int(fsize),
		StartTime: time.Now(),
	}

	return p, nil
}

func (s *ProgressStatus) Read(p []byte) (int, error) {
	n, err := s.Reader.Read(p)
	s.bytesCopied += n

	if s.Callback != nil {
		s.Callback(s.Path, s.bytesCopied, s.Size, s.StartTime)
	}

	return n, err
}

func getReaderSize(r io.Reader) (int64, error) {
	seeker, ok := r.(io.Seeker)
	if !ok {
		return 0, nil
	}

	// Save current position
	current, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// Seek to end to get size
	end, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Restore original position
	_, err = seeker.Seek(current, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return end, nil
}
