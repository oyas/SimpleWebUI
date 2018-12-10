package parser

import (
	"bufio"
	"gopkg.in/russross/blackfriday.v2"
	"os"
)

func Read(filename string, parse bool) string {
	var fp *os.File
	var err error
	var buf string

	fp, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		buf += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if parse {
		// pre parse in our parser
		mark := Parse(buf)
		// parse markdown with blackfriday
		unsafe := blackfriday.Run([]byte(mark))
		//html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
		return string(unsafe)
	}

	return buf
}
