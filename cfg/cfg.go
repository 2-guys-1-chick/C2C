package cfg

func GetPort() int {
	return 14975
}

func GetWsPort() int {
	return 8080
}

func GetVehicleId() string {
	return GetValue(VEHICLE_ID)
}
