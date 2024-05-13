package tools

import (
	"fmt"
	"time"
)

// GetTimeStr 获取时间字符串
func GetTimeStr(tm time.Time) string {
	return tm.Format("2006-01-02 15:04:05")
}

// Time json marsh 重写
type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	tmp := string(data)
	str := `"2006-01-02 15:04:05"`
	if len(tmp) <= 20 {
		str = `"2006-01-02"`
	}
	if tmp != `""` {
		now, err1 := time.ParseInLocation(str, tmp, time.Local)
		err = err1
		*t = Time{now}
	}

	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

func (t Time) String() string {
	return t.Format("2006-01-02 15:04:05")
}
