package servertype

import (
	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
)

type UpdateScooterStatusRequest struct {
	EventType string `json:"eventType" validate:"required,ascii"`
	ScooterID string `json:"scooterID" validate:"required,uuid4"`
	Latitude  string `json:"latitude"  validate:"required,latitude"`
	Longitude string `json:"longitude" validate:"required,longitude"`
	Timestamp string `json:"timestamp" validate:"required,datetime=2006-01-02T15:04:05Z"`
}

type FindScootersRequest struct {
	BottomLeft scootertype.Location `json:"bottomLeft" validate:"len=2,dive,required"`
	TopRight   scootertype.Location `json:"topRight" validate:"len=2,dive,required"`
}
