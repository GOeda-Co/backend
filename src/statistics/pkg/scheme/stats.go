package schemes

import "time"

type Interval struct {
	DtStart time.Time `json:"dt_start"`
	DtEnd time.Time `json:"dt_end"`
}