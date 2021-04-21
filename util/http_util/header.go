package http_util

import (
	"errors"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/http"
	"strings"
)

type Charset uint8

const (
	GB2312 Charset = iota
	GBK
	GB18030
	UTF8
)

var CharsetNames = map[string]Charset{
	"GB2312":GB2312,
	"GBK":GBK,
	"GB18030":GB18030,
	"UTF-8":UTF8,
	"UTF8":UTF8,
}

func GetCharset(head http.Header) Charset {
	if typeArr, ok := head["Content-Type"]; ok {
		for _, contentType := range typeArr {
			// 去掉空格
			contentType = strings.ReplaceAll(contentType, " ", "")
			contentType = strings.ToUpper(contentType)
			if idx := strings.LastIndex(contentType, "CHARSET="); idx > -1 {
				charset := contentType[idx+8:]
				if val, ok := CharsetNames[charset]; ok {
					return val
				}
			}
		}
	}
	return UTF8
}

/**
将其他编码的文本转换成utf8编码
 */
func Convert2UTF8(text []byte, charset Charset) ([]byte, error) {
	if charset == UTF8 {
		return text, nil
	}
	var encode encoding.Encoding = nil

	switch charset {
	case GB2312:
		encode = simplifiedchinese.HZGB2312
	case GBK:
		encode = simplifiedchinese.GBK
	case GB18030:
		encode = simplifiedchinese.GB18030
	}

	if encode != nil {
		if utf8Text, err := encode.NewDecoder().Bytes(text); err != nil {
			return nil, err
		} else {
			return utf8Text, nil
		}
	}

	return nil, errors.New("the charset is not supported")
}
