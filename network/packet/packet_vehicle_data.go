package packet

import (
	"bytes"
	"strconv"

	"fmt"

	"github.com/kellydunn/golang-geo"
)

type driveMode string

const (
	DriveModeAutopilot driveMode = "AUTO"
	DriveModeManual    driveMode = "MAN"
)

type VehicleData struct {
	Model           string    `json:"model"`
	ManufactureYear int       `json:"manufacture_year"`
	Speed           float64   `json:"speed"` // In Km/h
	Geo             geo.Point `json:"geo"`
	Weight          float64   `json:"weight"`
	TireWear        float64   `json:"tire_wear"`
	DriveMode       driveMode `json:"drive_mode"`
}

func (vd VehicleData) Encode() []byte {
	var bfr bytes.Buffer
	writeSubseparatedValue(&bfr, []byte(vd.Model))
	writeSubseparatedValue(&bfr, []byte(strconv.Itoa(vd.ManufactureYear)))
	writeSubseparatedValue(&bfr, formatFloat(vd.Speed, 3))
	writeSubseparatedValue(&bfr, gpsBytes(vd.Geo))
	writeSubseparatedValue(&bfr, formatFloat(vd.Weight, 4))
	writeSubseparatedValue(&bfr, formatFloat(vd.TireWear, 4))
	writeSubseparatedValue(&bfr, []byte(vd.DriveMode))

	return bfr.Bytes()
}

func (vd *VehicleData) Decode(bts []byte) error {
	parts := bytes.Split(bts, []byte{innerSubseparator})
	const mustPartsCount = 7
	if len(parts) != mustPartsCount {
		return fmt.Errorf("Vehicle Data: Unexpected number of parts, expected %d, received %d.", mustPartsCount, len(parts))
	}

	var err error
	vd.Model = string(parts[0])
	vd.ManufactureYear, err = strconv.Atoi(string(parts[1]))
	if err != nil {
		return err
	}

	vd.Speed, err = parseFloat(parts[2])
	if err != nil {
		return err
	}

	geoPtr, err := bytesGeo(parts[3])
	if err != nil {
		return err
	}

	vd.Geo = *geoPtr

	vd.Weight, err = parseFloat(parts[4])
	if err != nil {
		return err
	}

	vd.TireWear, err = parseFloat(parts[5])
	if err != nil {
		return err
	}

	vd.DriveMode = driveMode(parts[6])
	return nil
}

func bytesGeo(bts []byte) (*geo.Point, error) {
	parts := bytes.Split(bts, []byte{','})
	const mustPartsCount = 2
	if len(parts) != mustPartsCount {
		return nil, fmt.Errorf("Geo: Unexpected number of parts, expected %d, received %d.", mustPartsCount, len(parts))
	}

	lat, err := parseFloat(parts[0])
	if err != nil {
		return nil, err
	}

	lng, err := parseFloat(parts[1])
	if err != nil {
		return nil, err
	}

	return geo.NewPoint(lat, lng), nil
}

func gpsBytes(gps geo.Point) []byte {
	const prec = 7
	var bfr bytes.Buffer
	bfr.Write(formatFloat(gps.Lat(), prec))
	bfr.WriteByte(',')
	bfr.Write(formatFloat(gps.Lng(), prec))
	return bfr.Bytes()
}

func formatFloat(float float64, prec int) []byte {
	str := strconv.FormatFloat(float, 'f', prec, 64)
	return []byte(str)
}

func parseFloat(floatBts []byte) (float64, error) {
	return strconv.ParseFloat(string(floatBts), 64)
}
