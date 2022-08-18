package testutils

import (
	"testing"

	"github.com/pysf/special-umbrella/internal/scooter/scootertype"
)

const (
	RecBottomLeftLat float64 = 52.519121
	RecBottomLeftLng float64 = 13.402381
	RecTopRightLat   float64 = 52.538528
	RecTopRightLng   float64 = 13.427941
	BerlinCenterLat  float64 = 52.520008
	BerlinCenterLng  float64 = 13.404954
)

var BottomLeft = scootertype.Location{RecBottomLeftLat, RecBottomLeftLng}
var TopRight = scootertype.Location{RecTopRightLat, RecTopRightLng}

func IsInRectangle(bottomLeft, topRigth, coordinates scootertype.Location, t *testing.T) bool {
	lat := coordinates[0]
	if !(bottomLeft[0] < lat || lat < topRigth[0]) {
		t.Error("IsInRectangle() err= Returned location latitude is not in range")
		return false
	}

	lng := coordinates[1]
	if !(bottomLeft[1] < lng || lng < topRigth[1]) {
		t.Error("IsInRectangle() err= Returned location longitude is not in range")
		return false
	}
	return true
}
