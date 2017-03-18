package cfg

type configType uint

const (
	VEHICLE_ID configType = iota
)

var (
	configurations = map[configType]*spec{
		VEHICLE_ID: {
			"VEHICLE",
			"vehicle",
			"veh1",
		},
	}
)
