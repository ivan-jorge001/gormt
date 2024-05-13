package tools

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// AsString 转成string
func AsString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return GetTimeStr(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		{
			b, _ := json.Marshal(v)
			return string(b)
		}
	}
}

// DbcToSbc 全角转半角
func DbcToSbc(str string) string {
	numConv := unicode.SpecialCase{
		unicode.CaseRange{
			Lo: 0x3002, // Lo 全角句号
			Hi: 0x3002, // Hi 全角句号
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x002e - 0x3002, // LowerCase 转成半角句号
				0,               // TitleCase
			},
		},
		//
		unicode.CaseRange{
			Lo: 0xFF01, // 从全角！
			Hi: 0xFF19, // 到全角 9
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x0021 - 0xFF01, // LowerCase 转成半角
				0,               // TitleCase
			},
		},
		unicode.CaseRange{
			Lo: 0xff21, // Lo: 全角 Ａ
			Hi: 0xFF5A, // Hi:到全角 ｚ
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x0041 - 0xff21, // LowerCase 转成半角
				0,               // TitleCase
			},
		},
	}

	return strings.ToLowerSpecial(numConv, str)
}
