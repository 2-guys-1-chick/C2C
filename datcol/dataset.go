package datcol

import (
	"fmt"

	"github.com/2-guys-1-chick/c2c/network/packet"
	"github.com/kellydunn/golang-geo"
)

type pointTime struct {
	geo *geo.Point
	ms  int
}

var vehiclesMovementData = map[string][]pointTime{
	"veh1": {
		{geo: geo.NewPoint(50.079176, 14.432763), ms: 0},
		{geo: geo.NewPoint(50.083114, 14.434791), ms: 25000},
		{geo: geo.NewPoint(50.084566, 14.435928), ms: 12000},
		{geo: geo.NewPoint(50.086294, 14.437591), ms: 15000},
		{geo: geo.NewPoint(50.086294, 14.437591), ms: 10000},
		{geo: geo.NewPoint(50.08511, 14.440831), ms: 20000},
		{geo: geo.NewPoint(50.082033, 14.439511), ms: 20000},
		{geo: geo.NewPoint(50.078143, 14.436346), ms: 40000},
		{geo: geo.NewPoint(50.078143, 14.436346), ms: 10000},
		{geo: geo.NewPoint(50.077207, 14.440734), ms: 20000},
	},
	"veh2": {
		{geo: geo.NewPoint(50.078398, 14.442161), ms: 0},
		{geo: geo.NewPoint(50.078515, 14.445047), ms: 20000},
		{geo: geo.NewPoint(50.081468, 14.443996), ms: 22000},
		{geo: geo.NewPoint(50.081324, 14.444768), ms: 5000},
		{geo: geo.NewPoint(50.082735, 14.443942), ms: 10000},
		{geo: geo.NewPoint(50.084417, 14.44166), ms: 20000},
		{geo: geo.NewPoint(50.084487, 14.440517), ms: 6000},
		{geo: geo.NewPoint(50.085255, 14.440863), ms: 7000},
		{geo: geo.NewPoint(50.087272, 14.435091), ms: 30000},
		{geo: geo.NewPoint(50.087169, 14.428461), ms: 20000},
	},
}

var vehicleData = map[string]packet.VehicleData{
	"veh1": {
		Model:           "Škoda Kodiaq",
		ManufactureYear: 2017,
		Weight:          0.85,
		TireWear:        0.15,
		DriveMode:       packet.DriveModeManual,
	},
	"veh2": {
		Model:           "Tesla Model S",
		ManufactureYear: 2016,
		Weight:          0.75,
		TireWear:        0.10,
		DriveMode:       packet.DriveModeAutopilot,
	},
}

var driverData = map[string]packet.DriverData{
	"veh1": {
		Moods: packet.Moods(packet.MoodTired, packet.MoodInexperienced),
	},
	"veh2": {
		Moods: packet.Moods(packet.MoodDistracted, packet.MoodImpetuous),
	},
}

func calculateMovement(vehicleId string, ms int) (gps *geo.Point, speed float64) {
	movementData, ok := vehiclesMovementData[vehicleId]
	if !ok {
		return nil, 0
	}

	cumulativeMs := 0
	for i, mvm := range movementData {
		nextStepMs := cumulativeMs + mvm.ms
		if nextStepMs <= ms {
			cumulativeMs = nextStepMs
			continue
		}

		if i < 1 {
			continue
		}

		msInPeriod := ms - cumulativeMs
		var nextPoint *geo.Point

		currentPoint := movementData[i-1].geo
		nextPoint = mvm.geo
		latDiff := currentPoint.Lat() - nextPoint.Lat()
		lngDiff := currentPoint.Lng() - nextPoint.Lng()
		portion := float64(msInPeriod) / float64(mvm.ms)

		_ = fmt.Println
		pnt := geo.NewPoint(currentPoint.Lat()-latDiff*portion, currentPoint.Lng()-lngDiff*portion)

		totalDistance := currentPoint.GreatCircleDistance(nextPoint)
		speed := totalDistance / (float64(mvm.ms) / 1000 / 3600)
		return pnt, speed
	}

	return nil, 0
}
