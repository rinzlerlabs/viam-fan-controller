# viam-fan-controller modular resource

[![Go](https://github.com/viam-soleng/viam-fan-controller/actions/workflows/go.yml/badge.svg)](https://github.com/viam-soleng/viam-fan-controller/actions/workflows/go.yml)

A module to control a fan with feedback from temperature sensors connected to Viam.

This module provides two models for different kinds of fan controls: [PWM](#pwm-fan) and [On/Off](#onoff-fan).

## PWM Fan

A PWM fan allows for variable speed control. Fans typically require too much power to drive directly off a micro-controller pin, so PWM fans typically have 3 pins: a power pin, a ground pin, and a PWM signal pin. The signal pin energizes the coils of the fan with each pulse. The longer the pulses are on, the faster the fan spins.

### Build and run PWM fan

To use this module, follow the instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `viam-soleng:fan:pwm` model from the [`viam-fan-controller` module](https://app.viam.com/module/viam-soleng/viam-fan-controller).

### Configure your PWM fan

> [!NOTE]
> Before configuring your fan, you must [create a machine](https://docs.viam.com/manage/fleet/robots/#add-a-new-robot).

Navigate to the **Config** tab of your machine’s page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `sensor` type, then select the `fan:pwm` model.
Click **Add module**, then enter a name for your fan and click **Create**.

On the new component panel, copy and paste the following attribute template into your fan’s **Attributes** box:

```json
{
    "board_name": "<your board name>",
    "fan_pin": "<pin number>",
    "sensor_name": "<name of your temperature sensor>",
    "sensor_value_field": "<your temp sensor field key>",
    "temperature_table": {
        "0": 0,
        "30": 50,
        "50": 100
    }
}
```

Edit the values in the template as necessary, then click **Save config**.

> [!NOTE]
> For more information, see [Configure a Machine](https://docs.viam.com/manage/configuration/).

#### Attributes

The following attributes are available for `viam-soleng:fan:pwm` fans:

| Name | Type | Inclusion | Description |
| ---- | ---- | --------- | ----------- |
| board_name | string | **Required** | The `name` of the board that provides access to the GPIO pin to control the fan. |
| fan_pin | string | **Required** | The name of the GPIO pin on the board the fan is connected to. _Use the pin number, **not** the GPIO number_. |
| sensor_name | string | **Required** | The name of the sensor that provides the temperature feedback. |
| sensor_value_field | string | **Required** | The key name of the temperature in the sensor as returned by `Readings()`. |
| sensor_value_regex | string | Optional | A Regular Expression to parse the temperature out of the value returned by `Readings()`. This is only required if the value is a string and contains any characters not part of a valid floating point number. |
| temperature_table | map\[string\]float64| **Required** | A table that defines the temperature/fan speed values. |

> [!NOTE]
> The units of the `temperature_table` and the units of the temperature returned by the sensor must match.

Example configuration:

```json
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

### Build and run on/off fan

To use this module, follow the instructions to [add a module from the Viam Registry](https://docs.viam.com/registry/configure/#add-a-modular-resource-from-the-viam-registry) and select the `viam-soleng:fan:onoff` model from the [`viam-fan-controller` module](https://app.viam.com/module/viam-soleng/viam-fan-controller).

### Configure your on/off fan

> [!NOTE]
> Before configuring your fan, you must [create a machine](https://docs.viam.com/manage/fleet/robots/#add-a-new-robot).

Navigate to the **Config** tab of your machine’s page in [the Viam app](https://app.viam.com/).
Click on the **Components** subtab and click **Create component**.
Select the `sensor` type, then select the `fan:onoff` model.
Click **Add module**, then enter a name for your fan and click **Create**.

On the new component panel, copy and paste the following attribute template into your fan’s **Attributes** box:

```json
{
    "board_name": "<your board name>",
    "fan_pin": "<pin number>",
    "sensor_name": "<name of your temperature sensor>",
    "sensor_value_field": "<your temp sensor field key>",
    "on_temperature": 50,
    "off_temperature": 45,
    "on_delay": 5,
    "off_delay": 7
}
```

Edit the values in the template as necessary, then click **Save config**.

> [!NOTE]
> For more information, see [Configure a Machine](https://docs.viam.com/manage/configuration/).

The following attributes are available for `viam-soleng:fan:onoff` fans:

| Name | Type | Inclusion | Description |
| ---- | -----| --------- | ----------- |
| board_name | string | **Required** | The `name` of the board that provides access to the GPIO pin to control the fan. |
| fan_pin | string | **Required** | The name of the GPIO pin on the board the fan is connected to. _Use the pin number, **not** the GPIO number_. |
| sensor_name | string | **Required** | The `name` of the sensor that provides the temperature feedback. |
| sensor_value_field | string | **Required** | The key name of the temperature in the sensor as returned by `Readings()`. |
| sensor_value_regex | string | Optional | A Regular Expression to parse the temperature out of the value returned by `Readings()`. This is only required if the value is a string and contains any characters not part of a valid floating point number. |
| on_temperature | float64 | **Required** | The temperature at which to turn the fan on. |
| off_temperature | float64 | **Required** | The temperature at which to turn the fan off. |
| on_delay | int64 | Optional | The number of seconds to wait to turn the fan on after it was last turned off. This prevents turning the fan on/off too quickly. |
| off_delay | int64 | Optional | The number of seconds to wait to turn the fan off after it was last turned on. This prevents turning the fan on/off too quickly. |

> [!NOTE]
> The units of the on_temperature/off_temperature and the units of the temperature returned by the sensor must match.

Example configuration:

```json
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

## Local development

To use the `viam-fan-controller` module with a local install, clone this repository to your machine’s computer, navigate to the `viam-fan-controller` directory, and run:

```go
go build
```

On your robot’s page in the [Viam app](https://app.viam.com/), enter
the [module’s executable path](/registry/create/#prepare-the-module-for-execution), then click **Add module**.
The name must use only lowercase characters.
Then, click **Save config**.

## Next steps

- To test your fan, go to the [**Control** tab](https://docs.viam.com/manage/fleet/robots/#control).
- To write code against your fan, use one of the [available SDKs](https://docs.viam.com/program/).
- To view examples using a sensor component, explore [these tutorials](https://docs.viam.com/tutorials/).