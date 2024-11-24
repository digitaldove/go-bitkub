package bitkub

import "time"

// TimestampV2 represents time as the number of seconds since January 1, 1970, UTC.
type TimestampV2 int64

// TimestampV3 represents time as the number of milliseconds since January 1, 1970, UTC.
type TimestampV3 int64

func NewTimestamp(v time.Time) TimestampV3 {
	return TimestampV3(v.UnixMilli())
}

func (t TimestampV3) Time() time.Time {
	return time.UnixMilli(int64(t))
}

func (t TimestampV3) String() string {
	return t.Time().String()
}

func (t *TimestampV3) Set(v time.Time) {
	*t = TimestampV3(v.UnixMilli())
}

func (t TimestampV2) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t TimestampV2) String() string {
	return t.Time().String()
}

func (t *TimestampV2) Set(v time.Time) {
	*t = TimestampV2(v.Unix())
}
