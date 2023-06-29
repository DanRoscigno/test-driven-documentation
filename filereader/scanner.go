package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Section struct {
	Name string
	Body string
}

type Parser struct {
	Sections []*Section
	index    map[string]*Section
}

func (p *Parser) Section(name string) *Section {
	if p.index == nil {
		return nil
	}
	section, _ := p.index[name]
	return section
}

func (p *Parser) Parse(r io.Reader) {
	p.Sections = nil
	p.index = map[string]*Section{}

	scan := bufio.NewScanner(r)
	scan.Split(bufio.ScanLines)

	var (
		section    *Section
		buf        bytes.Buffer
		tokenStart = []byte(`-- section:`)
		tokenEnd   = []byte(`-- end`)
	)

	for scan.Scan() {
		line := scan.Bytes()
		if !bytes.HasPrefix(line, tokenStart) {
			continue
		}

		name := bytes.TrimSpace(line[len(tokenStart):])
		section = &Section{Name: string(name)}

		for scan.Scan() {
			line = scan.Bytes()
			if bytes.HasPrefix(line, tokenEnd) {
				section.Body = buf.String()
				p.Sections = append(p.Sections, section)
				p.index[section.Name] = section
				buf.Reset()
				break
			}
			buf.Write(line)
		}
	}
}

const content = `
# some comments
# comment
-- section: foo
this is the content for section foo
-- end
other stuff
-- section: bar
this is the content for section bar
-- end
# end stuff
`

func main() {
	r := strings.NewReader(content)

	var parser Parser
	parser.Parse(r)

	fmt.Println(parser.Section("bar"))

	for i, sec := range parser.Sections {
		fmt.Printf("#%d: %+v\n", i, sec)
	}
}
