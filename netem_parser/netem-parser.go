package netem_parser

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type NetemData struct {
	Total     int
	Dropped   int
	Reordered int
	Bytes     int
}

func Parse(text string) (*NetemData, error) {
	netemData := &NetemData{}
	scanner := bufio.NewScanner(strings.NewReader(text))

	for lineNo := 0; scanner.Scan(); lineNo++ {
		var match bool
		line := scanner.Text()

		switch lineNo {
		case 0:
			match, _ = regexp.MatchString("^qdisc netem.*$", line)
			if !match {
				return nil, errors.New("not a netem stat blob")
			}
		case 1:
			match, _ = regexp.MatchString("^ Sent.*$", line)
			if !match {
				return nil, errors.New("malformed netem stats")
			}
			var dummy int
			fmt.Sscanf(line, " Sent %d bytes %d pkt (dropped %d, overlimits %d requeues %d)",
				&netemData.Bytes, &netemData.Total, &netemData.Dropped, &dummy, &netemData.Reordered)
			return netemData, nil
		}
	}
	return nil, errors.New("bad netem stats")
}
