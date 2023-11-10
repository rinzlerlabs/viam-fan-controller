package utils

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"go.viam.com/rdk/logging"
)

func ParseCurrentTemperatureFromReadings(ctx context.Context, readings map[string]interface{}, sensorValueField string, sensorValueRegex *regexp.Regexp, logger logging.Logger) (float64, error) {

	// The numeric conversions are easy, but the string conversion is a little more complicated
	switch readings[sensorValueField].(type) {
	case float32:
		return float64(readings[sensorValueField].(float32)), nil
	case float64:
		return readings[sensorValueField].(float64), nil
	case int:
		return float64(readings[sensorValueField].(int)), nil
	case int32:
		return float64(readings[sensorValueField].(int32)), nil
	case int64:
		return float64(readings[sensorValueField].(int64)), nil
	case string:
		// First cast it to a string
		rawCurrentTemp := readings[sensorValueField].(string)
		if rawCurrentTemp == "" {
			logger.Errorf("Error reading sensor, field %s not found", sensorValueField)
			return 0, fmt.Errorf("error reading sensor, field %s not found", sensorValueField)
		}
		var currentTempString string
		if sensorValueRegex != nil {
			// Now try to use the regex to parse out the value
			currentTempString := sensorValueRegex.FindString(rawCurrentTemp)
			if currentTempString == "" {
				logger.Errorf("Error reading sensor, no match to regex in %s", currentTempString)
				return 0, fmt.Errorf("error reading sensor, no match to regex in %s", currentTempString)
			}
		} else {
			// If we don't have a regex, just use the whole string
			currentTempString = rawCurrentTemp
		}
		// Now convert it to a float and return it
		return strconv.ParseFloat(currentTempString, 64)
	default:
		logger.Errorf("Error reading sensor, field %s is unknown type", sensorValueField)
		return 0, fmt.Errorf("error reading sensor, field %s is unknown type", sensorValueField)
	}
}
