package date

import "time"

type TimeDate struct {
}

func New() DateInterface {
	time.Local, _ = time.LoadLocation("America/Sao_Paulo")
	return &TimeDate{}
}

func (t *TimeDate) Now() time.Time {
	return time.Now()
}

func (t *TimeDate) NowInRfc3339() string {
	return time.Now().Format(time.RFC3339)
}
