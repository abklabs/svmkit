package payload

import (
	"bytes"
	"io"
	"io/fs"
	"strings"
	"text/template"
)

type PayloadFile struct {
	Path   string
	Reader io.Reader
	Mode   fs.FileMode
}

type Payload struct {
	RootPath    string
	Files       []PayloadFile
	DefaultMode fs.FileMode
}

func (p *Payload) Add(f PayloadFile) {
	if f.Mode == 0 {
		if p.DefaultMode != 0 {
			f.Mode = p.DefaultMode
		} else {
			f.Mode = 0644
		}
	}

	p.Files = append(p.Files, f)
}

func (p *Payload) AddString(path string, body string) {
	p.AddReader(path, strings.NewReader(body))
}

func (p *Payload) AddReader(path string, reader io.Reader) {
	p.Add(PayloadFile{Path: path, Reader: reader})
}

func (p *Payload) NewBuffer(info PayloadFile, inb []byte) *bytes.Buffer {
	b := bytes.NewBuffer(inb)
	info.Reader = b
	p.Add(info)

	return b
}

func (p *Payload) NewWriter(info PayloadFile) io.Writer {
	return p.NewBuffer(info, nil)
}

func (p *Payload) AddTemplate(path string, tmpl *template.Template, data any) error {
	w := p.NewWriter(PayloadFile{Path: path})
	err := tmpl.Execute(w, data)

	return err
}
