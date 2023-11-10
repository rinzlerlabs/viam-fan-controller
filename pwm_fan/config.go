package pwm_fan

import "errors"

type CloudConfig struct {
	BoardName        string             `json:"board_name"`
	FanPin           string             `json:"fan_pin"`
	SensorName       string             `json:"sensor_name"`
	SensorValueField string             `json:"sensor_value_field"`
	SensorValueRegex string             `json:"sensor_value_regex"`
	TemperatureTable map[string]float64 `json:"temperature_table"`
}

func (conf *CloudConfig) Validate(path string) ([]string, error) {
	if conf.BoardName == "" {
		return nil, errors.New("board_name is required")
	}

	if conf.FanPin == "" {
		return nil, errors.New("pin is required")
	}

	if conf.SensorName == "" {
		return nil, errors.New("sensor_name is required")
	}

	if conf.SensorValueField == "" {
		return nil, errors.New("sensor_value is required")
	}

	if conf.TemperatureTable == nil {
		return nil, errors.New("temperature_table is required")
	}

	return nil, nil
}
