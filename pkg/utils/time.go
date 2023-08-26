package utils

import (
	"beastpark/meetinginvitationservice/pkg/log"
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	// timeFormat = "2006-01-02 15:04:05"
	timeFormat = "2006-01-02 15:04"
)

type StringTime struct {
	time.Time
}

func (t *StringTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	if err != nil {
		log.GetInstance().Infoln("time unmarshal fail ", err.Error())
		return
	}
	*t = StringTime{now}
	return
}

func (t *StringTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t.Time).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t StringTime) String() string {
	return time.Time(t.Time).Format(timeFormat)
}

func (ts StringTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(ts.Time)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan valueof time.Time
func (ts *StringTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*ts = StringTime{value}
		return nil
	}

	return fmt.Errorf("can not convert %v ", v)
}
