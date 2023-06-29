package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const data = `
-- section: foo
this is the content for section foo
-- end
-- section: bar
this is the content for section bar
-- end
-- section: dan
Dan the man
-- end
`

func main() {
	r := strings.NewReader(data)
	fmt.Println(ReadSections(r))
}

const SectionPrefix = "-- section:"
const SectionSuffix = "-- end"

func ReadSections(r io.Reader) (map[string]string, error) {
	ss := bufio.NewScanner(r)
	sections := make(map[string]string)
	for ss.Scan() {
		// parse the first line
		openline := ss.Text()
		if !strings.HasPrefix(openline, SectionPrefix) {
			continue

		}
		name := strings.TrimSpace(openline[len(SectionPrefix):])
		// parse the content
		var content strings.Builder
		for ss.Scan() {
			line := ss.Text()
			if strings.HasPrefix(line, SectionSuffix) {
				break
			}
			fmt.Fprintln(&content, line)
		}
		sections[name] = content.String()
	}
	return sections, ss.Err()
}
