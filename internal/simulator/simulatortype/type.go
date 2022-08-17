package simulatortype

import (
	"time"
)

type SimulatorResponse struct {
	Message   string
	Timestamp time.Time
	Err       error
}
