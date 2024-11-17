package runner

import (
	"io"
	"io/fs"
	"strings"
)

type PayloadFile struct {
	Path   string
	Reader io.Reader
	Mode   fs.FileMode
}

type Payload struct {
	RootPath string
	Files    []PayloadFile
}

func (p *Payload) Add(f PayloadFile) {
	p.Files = append(p.Files, f)
}

func (p *Payload) AddString(path string, body string) {
	p.AddReader(path, strings.NewReader(body))
}

func (p *Payload) AddReader(path string, reader io.Reader) {
	p.Add(PayloadFile{Path: path, Reader: reader, Mode: 0644})
}