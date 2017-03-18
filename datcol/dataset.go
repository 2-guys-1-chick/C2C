package datcol

import (
	"github.com/kellydunn/golang-geo"
	"fmt"
)

type pointTime struct {
	geo *geo.Point
	ms  int
}

var movement = []pointTime{
	{geo: geo.NewPoint(50.083114, 14.434791), ms: 0},
	{geo: geo.NewPoint(50.084566, 14.435928), ms: 12000},
	{geo: geo.NewPoint(50.086294, 14.437591), ms: 15000},
	{geo: geo.NewPoint(50.086294, 14.437591), ms: 10000},
	{geo: geo.NewPoint(50.08511, 14.440831), ms: 20000},
	{geo: geo.NewPoint(50.082033, 14.439511), ms: 20000},
	{geo: geo.NewPoint(50.078143, 14.436346), ms: 40000},
	{geo: geo.NewPoint(50.078143, 14.436346), ms: 10000},
	{geo: geo.NewPoint(50.077207, 14.440734), ms: 20000},
}

func calculateMovement(ms int) (gps *geo.Point, speed float64) {
	cumulativeMs := 0
	for i, mvm := range movement {
		nextStepMs := cumulativeMs + mvm.ms
		if nextStepMs <= ms {
			cumulativeMs = nextStepMs
			continue
		}

		msInPeriod := ms - cumulativeMs
		var nextPoint *geo.Point
		if len(movement) < i+1 {
			return nil, 0
		}

		currentPoint := movement[i-1].geo
		nextPoint = mvm.geo
		latDiff := currentPoint.Lat() - nextPoint.Lat()
		lngDiff := currentPoint.Lng() - nextPoint.Lng()
		portion := float64(msInPeriod) / float64(mvm.ms)

		_ = fmt.Println
		//fmt.Printf("From %v to %v\n", currentPoint, nextPoint)
		pnt := geo.NewPoint(currentPoint.Lat()-latDiff*portion, currentPoint.Lng()-lngDiff*portion)

		totalDistance := currentPoint.GreatCircleDistance(nextPoint)
		speed := totalDistance / (float64(mvm.ms) / 1000 / 3600)
		return pnt, speed
	}

	return nil, 0
}
