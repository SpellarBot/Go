// Convertor
// @Author: Golion
// @Date: 2017.5

package utils

import (
	"strconv"
	"net/url"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"fmt"
)

func URLEncode(str string) string {
	return url.QueryEscape(str)
}

func URLDecode(str string) string {
	if output, err := url.QueryUnescape(str); err != nil {
		return output
	} else {
		return ""
	}
}

func JSONEncode(jsObj interface{}) (string, error) {
	output, err := json.Marshal(jsObj)
	if err != nil {
		return "", fmt.Errorf("[JSONEncode] error=[%v]", err.Error())
	}
	return string(output), nil
}

func JSONDecode(jsStr string) (*simplejson.Json, error) {
	jsObj, err := simplejson.NewJson([]byte(jsStr))
	if err != nil {
		return nil, fmt.Errorf("[JSONDecode] error=[%v]", err.Error())
	}
	return jsObj, nil
}

func Atoi(x string) int {
	i, err := strconv.Atoi(x)
	if err != nil {
		return int(0)
	} else {
		return i
	}
}

func Atoi32(x string) int32 {
	i, err := strconv.ParseInt(x, 10, 32)
	if err != nil {
		return int32(0)
	} else {
		return int32(i)
	}
}

func Atoi64(x string) int64 {
	i, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return int64(0)
	} else {
		return int64(i)
	}
}

func Atof32(x string) float32 {
	f, err := strconv.ParseFloat(x, 32)
	if err != nil {
		return float32(0)
	} else {
		return float32(f)
	}
}

func Atof64(x string) float64 {
	f, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return float64(0)
	} else {
		return float64(f)
	}
}