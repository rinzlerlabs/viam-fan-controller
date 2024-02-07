package on_off_fan

import "errors"

type CloudConfig struct {
	BoardName        string  `json:"board_name"`
	FanPin           string  `json:"fan_pin"`
	SensorName       string  `json:"sensor_name"`
	SensorValueKey   string  `json:"sensor_value_key"`
	SensorValueRegex string  `json:"sensor_value_regex"`
	OnTemperature    float64 `json:"on_temperature"`
	OffTemperature   float64 `json:"off_temperature"`
	OnDelay          int64   `json:"on_delay"`
	OffDelay         int64   `json:"off_delay"`
}

func (conf *CloudConfig) Validate(path string) ([]string, error) {
	if conf.BoardName == "" {
		return nil, errors.New("board_name is required")
	}

	if conf.FanPin == "" {
		return nil, errors.New("fan_pin is required")
	}

	if conf.SensorName == "" {
		return nil, errors.New("sensor_name is required")
	}

	if conf.SensorValueKey == "" {
		return nil, errors.New("sensor_value_key is required")
	}

	if conf.OnTemperature == 0 {
		return nil, errors.New("on_temperature is required")
	}

	if conf.OffTemperature == 0 {
		return nil, errors.New("off_temperature is required")
	}

	return nil, nil
}
