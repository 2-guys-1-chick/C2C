package datcol

import (
	"testing"
	"fmt"
	"time"
)

func TestCalculateMovement(t *testing.T) {
	printInTime(1000)
	printInTime(2000)
	printInTime(3000)
	printInTime(8000)
	printInTime(61000)
}

func TestCalculatePath(t *testing.T) {
	start := time.Now()
	for {
		diff := time.Now().Sub(start)
		pnt , speed := calculateMovement(int(diff.Seconds() * 1000))
		if pnt == nil {
			break
		}

		fmt.Println(pnt, speed)
		time.Sleep(time.Second)
	}
}


func printInTime(ms int) {
	gps, speed := calculateMovement(ms)
	fmt.Println(gps, speed)
}
