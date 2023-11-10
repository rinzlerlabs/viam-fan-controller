# viam-fan-controller

[![Go](https://github.com/viam-soleng/viam-fan-controller/actions/workflows/go.yml/badge.svg)](https://github.com/viam-soleng/viam-fan-controller/actions/workflows/go.yml)

A module to control fans based on temperature sensors connected to Viam.

This module offers 2 kinds of fan controls: [PWM](#pwm) and [On/Off](#onoff-fan).

## PWM

A PWM fan allows for variable speed control. Fans typically require too much power to drive directly off a micro-controller pin, so PWM fans typically have 3 pins, a power pin, a ground pin, and a signal pin. The signal pin energizes the coils of the fan with each pulse. The longer the pulses are on, the faster the fan spins.

The config for a PWM fan requires:
|Attribute|Required|Type|Description|
|---------|--------|----|-----------|
|board_name|Y|string|The name of the board that provides access to the GPIO pin to control the fan|
|fan_pin|Y|string|The name of the GPIO pin on the board the fan is connected to|
|sensor_name|Y|string|The name of the sensor that provides the temperature feedback|
|sensor_value_field|Y|string|The key name of the temperature in the sensor as returned by Readings()|
|sensor_value_regex|N|string|A Regular Expression to parse the temperature out of the value returned by Readings(). This is only required if the value is a string and contains any characters not part of a valid floating point number|
|temperature_table|Y|map\[string\]float64|A table that defines the temperature/fan speed values|

_Note: The units of the temperature_table and the units of the temperature returned by the sensor must match._

Example Config:
```
{
    "board_name": "pi",
    "fan_pin": "15",
    "sensor_name": "board_temps",
    "sensor_value_field": "soc_temp",
    "temperature_table": {
        "0": 0,
        "30": 50,
        "50": 100
    }
}
```

In this config, there is a sensor already configured with the name `board_temps` that is providing a field `soc_temp` returned in `Readings()`. Note that there is no `sensor_value_regex` because this sensor already returns a `float64` for the temperature.

## On/Off Fan

A simple on/off fan does just that, it is either on or off. This is a useful for driving larger fans that have their own external speed controllers or require more power than a micro-controller can provide. In cases like that, the GPIO pin will just drive a relay or a simple signal into the external motor controller.

The config for a simple On/Off fan requires:
|Attribute|Required|Type|Description|
|---------|--------|----|-----------|
|board_name|Y|string|The name of the board that provides access to the GPIO pin to control the fan|
|fan_pin|Y|string|The name of the GPIO pin on the board the fan is connected to|
|sensor_name|Y|string|The name of the sensor that provides the temperature feedback|
|sensor_value_field|Y|string|The key name of the temperature in the sensor as returned by Readings()|
|sensor_value_regex|N|string|A Regular Expression to parse the temperature out of the value returned by Readings(). This is only required if the value is a string and contains any characters not part of a valid floating point number|
|on_temperature|Y|float64|The temperature at which to turn the fan on|
|off_temperature|Y|float64|The temperature at which to turn the fan off|
|on_delay|N|int64|The number of seconds to wait to turn the fan on after it was last turned off. This prevents flapping of the fan on/off too quickly|
|off_delay|N|int64|The number of seconds to wait to turn the fan off after it was last turned on. This prevents flapping of the fan on/off too quickly|

_Note: The units of the on_temperature/off_temperature and the units of the temperature returned by the sensor must match._

Example Config:
```
{
    "board_name": "pi",
    "fan_pin": "15",
    "sensor_name": "board_temps",
    "sensor_value_field": "soc_temp",
    "on_temperature": 50,
    "off_temperature": 45,
    "on_delay": 5
}
```

In this config, there is a sensor already configured with the name `board_temps` that is providing a field `soc_temp` returned in `Readings()`. The fan will turn on when the `soc_temp` goes above 50 and will turn off again when the temperature goes below 45. After `soc_temp` exceeds 50, if the fan had previously been turned off less than 5 seconds ago, the fan will not turn on until 5 seconds has elapsed since the fan was turned off.
