package bitkub

import "time"

// Timestamp represents time as the number of seconds since January 1, 1970, UTC.
type Timestamp int64

func NewTimestamp(v time.Time) Timestamp {
	return Timestamp(v.Unix())
}

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

func (t Timestamp) String() string {
	return t.Time().String()
}

func (t *Timestamp) Set(v time.Time) {
	*t = Timestamp(v.Unix())
}
