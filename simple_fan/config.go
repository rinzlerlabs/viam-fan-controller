package simple_fan

type CloudConfig struct {
	BoardName        string  `json:"board_name"`
	Pin              string  `json:"pin"`
	SensorName       string  `json:"sensor_name"`
	SensorValue      string  `json:"sensor_value"`
	SensorValueRegex string  `json:"sensor_value_regex"`
	SensorUnits      string  `json:"sensor_units"`
	OnTemperature    float64 `json:"on_temperature"`
	OffTemperature   float64 `json:"off_temperature"`
	OnDelay          int64   `json:"on_delay"`
	OffDelay         int64   `json:"off_delay"`
}

func (conf *CloudConfig) Validate(path string) ([]string, error) {
	return nil, nil
}
