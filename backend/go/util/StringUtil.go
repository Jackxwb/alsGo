package util

import (
	bytes2 "bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"regexp"

	"io"
)

func TextGbkToUtf8(byteData []byte) ([]byte, error) {
	reader := bytes2.NewReader(byteData)
	charset := simplifiedchinese.GBK
	utf8Read := transform.NewReader(reader, charset.NewDecoder())
	all, err := io.ReadAll(utf8Read)
	return all, err
}

func RegexpFindStringGroup(temple, data string) map[string]string {
	comp := regexp.MustCompile(temple)
	submatch := comp.FindStringSubmatch(data)
	group := comp.SubexpNames()

	result := make(map[string]string)
	for i := 1; i < len(submatch); i++ {
		result[group[i]] = submatch[i]
	}
	return result
}

func RegexpFindString(temple, data string) bool {
	comp := regexp.MustCompile(temple)
	submatch := comp.FindStringSubmatch(data)
	if len(submatch) > 0 {
		return true
	}
	return false
}
